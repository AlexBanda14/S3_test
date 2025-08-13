package app

import (
	"S3Work/internal/config"
	"S3Work/internal/infastructure/s3"
	"S3Work/internal/usecase"
	"S3Work/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, lg *logger.ZapLogger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	clientS3, err := s3.NewS3Client(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.BucketName, ctx)
	if err != nil {
		lg.Logger.Errorw("NewS3Client err", "error", err)
		return fmt.Errorf("NewS3Client err: %v", err)
	}

	service := usecase.NewImagesService(clientS3, lg)

	err = service.UploadFolder(cfg.PathUpload, ctx)
	if err != nil {
		lg.Logger.Errorw("UploadFolder err", "error", err)
		return fmt.Errorf("UploadFolder err: %v", err)
	}
	lg.Logger.Infow("UploadFolder success", "path", cfg.PathUpload)

	err = service.DownloadAll(cfg.PathDownload, ctx)
	if err != nil {
		lg.Logger.Errorw("DownloadAll err", "error", err)
		return fmt.Errorf("DownloadAll err: %v", err)
	}
	lg.Logger.Infow("DownloadAll success", "path", cfg.PathDownload)

	<-sigChan
	cancel()
	lg.Logger.Infow("Shutting down application")
	return nil
}
