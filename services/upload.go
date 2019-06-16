package services

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/pocockn/awswrappers/s3"
	"github.com/pocockn/shouts-api/config"
	"mime/multipart"
	"net/http"
)

type (
	// Upload holds stuff related to uploading images to S3.
	Upload struct {
		config   config.Config
		s3Client *s3.Client
	}
)

// NewUpload creates and returns a new upload struct.
func NewUpload(config config.Config, s3Client *s3.Client) Upload {
	return Upload{
		config:   config,
		s3Client: s3Client,
	}
}

// UploadToS3 takes a file from a form and uploads it to S3.
func (u Upload) UploadToS3(sourceFile *multipart.FileHeader) error {
	key := sourceFile.Filename
	object := s3.NewObject(
		u.config.S3.Bucket,
		key,
		u.s3Client,
	)
	buffer := make([]byte, sourceFile.Size)

	src, err := sourceFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	_, err = src.Read(buffer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = object.Put(bytes.NewReader(buffer), http.DetectContentType(buffer))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}
