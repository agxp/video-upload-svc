package main

import (
	"log"
	"os"

	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"github.com/micro/go-micro"
	"time"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/agxp/cloudflix/tracer"
)

func main() {
	log.SetOutput(os.Stdout)

	zipkin_addr := os.Getenv("ZIPKIN_ADDR")
	hostname, _ := os.Hostname()
	InitTracer(zipkin_addr, hostname, "video_upload")

	// Creates a database connection and handles
	// closing it again before exit.
	s3, err := ConnectToS3()
	if err != nil {
		log.Fatalf("Could not connect to store: %v", err)
	}

	pg, err := ConnectToPostgres()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	repo := &UploadRepository{s3, pg}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("video_upload"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(trace.ServerWrapper),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Will comment this out now to save having to run this locally
	// publisher := micro.NewPublisher("user.created", srv.Client())

	// Register handler
	pb.RegisterUploadHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitTracer(zipkinURL string, hostPost string, serviceName string) {
	collector, err := zipkin.NewHTTPCollector(zipkinURL)
	if err != nil {
		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
		return
	}
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, hostPost, serviceName),
	)
	if err != nil {
		log.Fatalf("unable to create Zipkin tracer: %v", err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	return
}
