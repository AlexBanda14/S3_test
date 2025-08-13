package usecase

import (
	"S3Work/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"os"
	"path/filepath"
)

type ImagesInterface interface {
	UploadFolder(pathUpload string, ctx context.Context) error
	DownloadAll(pathDownload string, ctx context.Context) error
}

type ImagesService struct {
	storage StorageImage
	log     *logger.ZapLogger
}

func NewImagesService(storage StorageImage, log *logger.ZapLogger) *ImagesService {
	return &ImagesService{
		storage: storage,
		log:     log,
	}
}

func (s *ImagesService) UploadFolder(pathUpload string, ctx context.Context) error {
	files, err := os.ReadDir(pathUpload)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(pathUpload, file.Name())
		img, err := imaging.Open(filePath)
		if err != nil {
			s.log.Logger.Errorw("could not open image", "path", filePath, "err", err)
			continue
		}

		resize := imaging.Resize(img, 800, 0, imaging.Lanczos)

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, resize, &jpeg.Options{Quality: 80})
		if err != nil {
			s.log.Logger.Errorw("could not encode image", "path", filePath, "err", err)
			continue
		}

		err = s.storage.PutObject(file.Name(), buf, ctx)
		if err != nil {
			s.log.Logger.Errorw("could not put image", "name", file.Name(), "path", filePath, "err", err)
			continue
		}
		s.log.Logger.Infow("file uploaded successfully", "name", file.Name(), "path", filePath)
	}
	return nil
}

func (s *ImagesService) DownloadAll(pathDownload string, ctx context.Context) error {
	err := os.MkdirAll(pathDownload, os.ModePerm)
	if err != nil {
		s.log.Logger.Errorw("could not create directory", "path", pathDownload, "err", err)
		return fmt.Errorf("could not create directory: %w", err)
	}

	keys, err := s.storage.ListObject(ctx)
	if err != nil {
		s.log.Logger.Errorw("could not list objects", "err", err)
		return fmt.Errorf("could not list objects: %w", err)
	}

	for _, key := range keys {

		data, err := s.storage.GetObject(key, ctx)
		if err != nil {
			s.log.Logger.Errorw("could not get object", "key", key, "err", err)
			continue
		}

		filePath := filepath.Join(pathDownload, key)
		err = os.WriteFile(filePath, data.Bytes(), 0644)
		if err != nil {
			s.log.Logger.Errorw("could not write file", "path", filePath, "err", err)
			continue
		}
		s.log.Logger.Infow("file download successfully", "name", key, "path", filePath)
	}
	return nil
}
