package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"golang.org/x/net/context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type service struct {
	repo Repository
	tracer *opentracing.Tracer
	logger *zap.Logger
}

func (srv *service) S3Request(ctx context.Context, req *pb.Request, res *pb.Response) error {
	sp, _ := opentracing.StartSpanFromContext(context.TODO(), "S3Request_Service")

	logger.Info("Request for S3Request_Service received")
	defer sp.Finish()

	url, err := srv.repo.S3Request(sp.Context(), req.Filename)
	if err != nil {
		logger.Error("failed S3Request", zap.Error(err))
		return err
	}

	res.PresignedUrl = url
	return nil
}

func (srv *service) UploadFile(ctx context.Context, req *pb.UploadRequest, res *pb.UploadResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.TODO(), "UploadFile_Service")
	logger.Info("Request for UploadFile_Service received")
	defer sp.Finish()

	id, filePath, err := srv.repo.WriteVideoProperties(sp.Context(), req.Filename, req.Title, req.Description)
	if err != nil {
		logger.Error("failed WriteVideoProperties", zap.Error(err))
		return err
	}

	url, err := srv.repo.S3Request(sp.Context(), filePath)
	if err != nil {
		logger.Error("failed S3Request", zap.Error(err))
		return err
	}

	res.Id = id
	res.PresignedUrl = url

	return nil
}


func (srv *service) WriteVideoProperties(ctx context.Context, req *pb.PropertyRequest, res *pb.PropertyResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.TODO(), "WriteVideoProperties_Service")
	defer sp.Finish()
	return nil
}

func (srv *service) UploadFinish(ctx context.Context, req *pb.UploadFinishRequest, res *pb.UploadFinishResponse) error {
	sp, _ := opentracing.StartSpanFromContext(context.TODO(),"UploadFile_Service")
	logger.Info("Request for UploadFinish_Service received")

	defer sp.Finish()

	err := srv.repo.UploadFinish(sp.Context(), req.Id)
	if err != nil {
		logger.Error("failed UploadFinish", zap.Error(err))
		return err
	}
	return nil
}