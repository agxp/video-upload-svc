package main

import (
	"log"
	"os"

	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"github.com/micro/go-micro"
	"time"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	micro_opentracing "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
)

var (
	tracer *opentracing.Tracer
)


func main() {

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return
	}

	t, closer, err := cfg.NewTracer()
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	tracer = &t
	opentracing.SetGlobalTracer(t)
	defer closer.Close()

	(*tracer).StartSpan("init_tracing").Finish()
	// continue main()

	log.SetOutput(os.Stdout)

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

	repo := &UploadRepository{s3, pg, tracer}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("video_upload"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(micro_opentracing.NewHandlerWrapper(*tracer)),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Will comment this out now to save having to run this locally
	// publisher := micro.NewPublisher("user.created", srv.Client())

	// Register handler
	pb.RegisterUploadHandler(srv.Server(), &service{repo, tracer})

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}