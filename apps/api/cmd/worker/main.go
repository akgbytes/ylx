package main

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity"
	"github.com/akgbytes/ylx/internal/platform/config"
	"github.com/akgbytes/ylx/internal/platform/logger"
	"github.com/akgbytes/ylx/internal/platform/mailer"
	"github.com/akgbytes/ylx/internal/platform/redis"
)

const (
	workerConcurrency = 5
	emailQueueWeight  = 2
)

func main() {
	bootstrapLogger := logger.BootstrapLogger()

	cfg, err := config.Load()
	if err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("bootstrap worker")
	}

	log, err := logger.New(&cfg.Log)
	if err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("bootstrap worker")
	}

	if err := run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("worker exited")
	}
}

func run(cfg *config.Config, log zerolog.Logger) error {
	redisCtx, redisCancel := context.WithTimeout(context.Background(), cfg.Redis.ConnectTimeout)

	rdb, err := redis.NewClient(redisCtx, cfg.Redis)
	redisCancel()

	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			log.Err(err).Msg("close redis")
		}
	}()

	identityWorker := identity.NewWorker(identity.WorkerDeps{
		Config: cfg,
		Redis:  rdb,
		Sender: mailer.NewResendSender(cfg.Email.ResendAPIKey, cfg.Email.From),
		Logger: log,
	})

	mux := asynq.NewServeMux()
	identityWorker.RegisterTasks(mux)

	server := asynq.NewServerFromRedisClient(rdb, asynq.Config{
		Concurrency: workerConcurrency,
		Queues: map[string]int{
			identity.Queue: emailQueueWeight,
		},
	})

	log.Info().Msg("worker started")

	if err := server.Run(mux); err != nil {
		return fmt.Errorf("run worker: %w", err)
	}

	return nil
}
