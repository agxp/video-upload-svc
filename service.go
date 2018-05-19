package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"golang.org/x/net/context"
	"log"
	"os"
	"github.com/opentracing/opentracing-go"
)

type service struct {
	repo Repository
	tracer *opentracing.Tracer
}

func (srv *service) S3Request(ctx context.Context, req *pb.Request, res *pb.Response) error {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "S3Request_Service")
	defer sp.Finish()

	log.SetOutput(os.Stdout)

	url, err := srv.repo.S3Request(sp, req.Filename)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	res.PresignedUrl = url
	log.Print("res", res.PresignedUrl)

	return nil
}

func (srv *service) UploadFile(ctx context.Context, req *pb.UploadRequest, res *pb.UploadResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "UploadFile_Service")
	defer sp.Finish()
	log.SetOutput(os.Stdout)

	id, filePath, err := srv.repo.WriteVideoProperties(sp, req.Filename, req.Title, req.Description)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url, err := srv.repo.S3Request(sp, filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	res.Id = id
	res.PresignedUrl = url
	log.Print("res", res)

	return nil
}


func (srv *service) WriteVideoProperties(ctx context.Context, req *pb.PropertyRequest, res *pb.PropertyResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "WriteVideoProperties_Service")
	defer sp.Finish()
	return nil
}

func (srv *service) UploadFinish(ctx context.Context, req *pb.UploadFinishRequest, res *pb.UploadFinishResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.Background(),"UploadFile_Service")
	defer sp.Finish()
	err := srv.repo.UploadFinish(sp, req.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}