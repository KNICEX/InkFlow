package service

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"io"
)

type FileService interface {
	Upload(ctx context.Context, reader io.Reader) (string, error)
}

type CloudinaryFileService struct {
	client *cloudinary.Cloudinary
}

func (c *CloudinaryFileService) Upload(ctx context.Context, reader io.Reader) (string, error) {
	res, err := c.client.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return res.URL, nil
}
