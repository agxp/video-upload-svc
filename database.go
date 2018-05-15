package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
	"fmt"
)

func ConnectToS3() (*minio.Client, error) {
	log.SetOutput(os.Stdout)
	endpoint := os.Getenv("MINIO_URL")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup
	log.Printf("%#v\n", "hello, i just finished printing minioclient")

	buckets, err := minioClient.ListBuckets()
	if err != nil {
		fmt.Println(err)
	}
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
	log.Printf("%#v\n", "hello, i just finished printing all the buckets")

	return minioClient, nil
}