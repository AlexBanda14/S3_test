package main

import (
	"S3Work/internal/app"
	"S3Work/internal/config"
	"S3Work/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	lg, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer lg.Close()

	err = app.Run(cfg, lg)
	if err != nil {
		lg.Logger.Errorw("error in Run()", "error", err)
	}
}
