package main

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/email"
	"github.com/akgbytes/ylx/internal/logger"
	"github.com/akgbytes/ylx/internal/redis"
	"github.com/akgbytes/ylx/internal/tasks"
	"github.com/akgbytes/ylx/internal/worker"
)

const (
	emailWorkerConcurrency = 5
	emailQueueWeight       = 1
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
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.RedisConnectTimeout)
	rdb, err := redis.NewClient(ctx, cfg.Redis)
	cancel()
	if err != nil {
		return fmt.Errorf("connect Redis: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error().Err(err).Msg("close Redis")
		}
	}()

	emailSender := email.NewResendSender(cfg.Email.ResendAPIKey, cfg.Email.From)
	emailWorker := worker.NewEmailWorker(emailSender, rdb)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSendSignUpOTP, emailWorker.HandleSendSignUpOTP)

	workerServer := asynq.NewServerFromRedisClient(rdb, asynq.Config{
		Concurrency: emailWorkerConcurrency,
		Queues: map[string]int{
			"email": emailQueueWeight,
		},
	})

	log.Info().Msg("email worker started")

	if err := workerServer.Run(mux); err != nil {
		return fmt.Errorf("run email worker: %w", err)
	}

	return nil
}
