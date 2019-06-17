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

// UploadToS3 takes a number of files from a form and uploads it to S3.
func (u Upload) MultiUpload(files ...*multipart.FileHeader) error {
	errorChannel := make(chan error)

	for _, f := range files {
		go func(f *multipart.FileHeader) {
			key := f.Filename
			object := s3.NewObject(
				u.config.S3.Bucket,
				key,
				u.s3Client,
			)
			buffer := make([]byte, f.Size)

			src, err := f.Open()
			if err != nil {
				errorChannel <- echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			_, err = src.Read(buffer)
			if err != nil {
				errorChannel <- echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			err = object.Put(bytes.NewReader(buffer), http.DetectContentType(buffer))
			if err != nil {
				errorChannel <- echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			errorChannel <- nil
		}(f)
	}

	for range files {
		if err := <-errorChannel; err != nil {
			return err
		}
	}

	return nil
}
