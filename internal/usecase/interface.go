package usecase

import (
	"bytes"
	"context"
)

type StorageImage interface {
	PutObject(key string, data *bytes.Buffer, ctx context.Context) error
	GetObject(key string, ctx context.Context) (*bytes.Buffer, error)
	ListObject(ctx context.Context) ([]string, error)
}
