package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"golang.org/x/net/context"
	"log"
	"os"
)

const topic = "user.created"

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
	res = url

	return nil
}
