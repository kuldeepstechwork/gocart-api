package services_test

import (
	"mime/multipart"
	"testing"

	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"github.com/kuldeepstechwork/gocart-api/test/mocks"
	"go.uber.org/mock/gomock"
)

func TestUploadService_UploadProductImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockUploadProvider(ctrl)
	s := services.NewUploadService(mockProvider)

	t.Run("Success", func(t *testing.T) {
		file := &multipart.FileHeader{Filename: "test.jpg"}
		mockProvider.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return("http://example.com/test.jpg", nil)

		url, err := s.UploadProductImage(1, file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "http://example.com/test.jpg" {
			t.Errorf("expected URL http://example.com/test.jpg, got %s", url)
		}
	})

	t.Run("InvalidExtension", func(t *testing.T) {
		file := &multipart.FileHeader{Filename: "test.txt"}
		_, err := s.UploadProductImage(1, file)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
