package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
	"time"
	"crypto/md5"
)

type Repository interface {
	S3Request(filename string) (string, error)
	UploadFile(filename string) (string, error)
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

func (repo *UploadRepository) UploadFile(filename string) (string, error) {
	log.SetOutput(os.Stdout)
	log.Printf("%#v\n", "filename: " + filename)
	
	objectName := time.Now().String() + "_" + filename
	filePath := md5.Sum(objectName) "/" + filename

	n, err := repo.s3.FPutObject("videos", objectName, filePath)
	if err != nil {
		log.Printf("%#v\n", "FAILED: with filename: " + filename)
		log.Print(objectName)
		log.Print(filePath)
		log.Fatal(err)
		return err.Error(), err
	}

	log.Printf("Uploaded %s of size %d\n", objectName, n)

	return "", nil
}