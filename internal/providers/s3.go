package providers

import (
	"context"

	"github.com/kuldeepstechwork/gocart-api/internal/config"
)

// S3Provider is a placeholder for an S3 upload provider.
type S3Provider struct{}

// NewS3Provider creates a new S3 provider.
func NewS3Provider(cfg *config.Config) *S3Provider {
	return &S3Provider{}
}

// Upload uploads a file to S3.
func (p *S3Provider) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	return "", nil
}
