package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"golang.org/x/net/context"
	"log"
	"os"
)

type service struct {
	repo Repository
}

func (srv *service) S3Request(ctx context.Context, req *pb.Request, res *pb.Response) error {
	log.SetOutput(os.Stdout)

	url, err := srv.repo.S3Request(req.Filename)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	res.PresignedUrl = url
	log.Print("res", res.PresignedUrl)

	return nil
}

func (srv *service) UploadFile(ctx context.Context, req *pb.Request, res *pb.Response) error {
	log.SetOutput(os.Stdout)

	error, err := srv.repo.UploadFile(req.Filename)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	res.Error = error
	log.Print("res", res.error)

	return nil
}
