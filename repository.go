package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"github.com/minio/minio-go"
	"log"
	"os"
	"time"
)

type Repository interface {
	S3Request(filename string) (*pb.Response, error)
}

type UploadRepository struct {
	s3 *minio.Client
}

func (repo *UploadRepository) S3Request(filename string) (*pb.Response, error) {
	log.SetOutput(os.Stdout)
	var res *pb.Response
	log.Println("filename: ", filename)
	presignedURL, err := repo.s3.PresignedPutObject("videos", filename, time.Duration(1000)*time.Second)
	if err != nil {
		log.Fatalln(filename)
		return nil, err
	}

	res.PresignedUrl = presignedURL.String()

	return res, nil
}
