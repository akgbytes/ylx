package main

import (
	"github.com/akgbytes/ylx/internal/app"
	"github.com/akgbytes/ylx/internal/bootstrap"
)

func main() {
	bootstrapLogger := bootstrap.NewBootstrapLogger()

	runtime, err := bootstrap.Load()
	if err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("bootstrap application")
	}

	app := app.NewApplication(runtime.Config, runtime.Logger)

	if err := app.Run(); err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("server exited")
	}
}
