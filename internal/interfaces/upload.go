package interfaces

import "context"

// UploadProvider is an interface for uploading files.
type UploadProvider interface {
	Upload(ctx context.Context, filename string, data []byte) (string, error)
}
