package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
)

func ConnectToS3() (*minio.Client, error) {
	log.SetOutput(os.Stdout)
	endpoint := os.Getenv("MINIO_URL")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup
	return minioClient, nil
}