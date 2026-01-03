package providers

import (
	"context"
)

// LocalUploadProvider is a placeholder for a local upload provider.
type LocalUploadProvider struct {
	path string
}

// NewLocalUploadProvider creates a new local upload provider.
func NewLocalUploadProvider(path string) *LocalUploadProvider {
	return &LocalUploadProvider{path: path}
}

// Upload uploads a file to the local filesystem.
func (p *LocalUploadProvider) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	return "", nil
}
