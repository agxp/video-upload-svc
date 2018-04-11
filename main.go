package main

import (
	"log"
	"os"

	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	log.SetOutput(os.Stdout)

	// Creates a database connection and handles
	// closing it again before exit.
	s3, err := ConnectToS3()
	if err != nil {
		log.Fatalf("Could not connect to store: %v", err)
	}

	repo := &UploadRepository{s3}

	// Create a new service. Optionally include some options here.
	srv := k8s.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("cloudflix.api.video_upload"),
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
