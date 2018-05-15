package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
	"time"
)

type Repository interface {
	S3Request(filename string) (string, error)
}

type UploadRepository struct {
	s3 *minio.Client
}

// should return string
func (repo *UploadRepository) S3Request(filename string) (string, error) {
	log.SetOutput(os.Stdout)
	log.Printf("%#v\n", "filename: " + filename)

	presignedURL, err := repo.s3.PresignedPutObject("videos", filename, time.Duration(1000)*time.Second)
	if err != nil {
		log.Printf("%#v\n", "FAILED: with filename: " + filename)
		log.Fatal(err)
		return "", err
	}

	log.Print(presignedURL)

	return presignedURL.String(), nil
}
