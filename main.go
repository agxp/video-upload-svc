package main

import (
	pb "github.com/agxp/cloudflix/video-upload-svc/proto"
	"github.com/micro/go-micro"
	"time"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	micro_opentracing "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	zapWrapper "github.com/uber/jaeger-client-go/log/zap"
	"github.com/agxp/cloudflix/video-encoding-svc/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/broker"
	"log"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	tracer *opentracing.Tracer
	logger         *zap.Logger
	metricsFactory metrics.Factory
)


func main() {

	logger, _ = zap.NewDevelopment()
	metricsFactory = prometheus.New()

	zapLogger := logger.With(zap.String("service", "video-upload-svc"))
	jeagerLogger := zapWrapper.NewLogger(zapLogger)

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		zapLogger.Error("Could not parse Jaeger env vars: %s", zap.Error(err))
		return
	}

	t, closer, err := cfg.NewTracer(
		jaegercfg.Metrics(metricsFactory),
		jaegercfg.Logger(jeagerLogger),
	)
	if err != nil {
		jeagerLogger.Infof("Could not initialize jaeger tracer: %s", err)
		return
	}

	tracer = &t
	opentracing.SetGlobalTracer(t)
	defer closer.Close()

	(*tracer).StartSpan("init_tracing").Finish()
	// continue main()

	// Creates a database connection and handles
	// closing it again before exit.
	s3, err := ConnectToS3()
	if err != nil {
		jeagerLogger.Error("Could not connect to store: " + err.Error())
	}

	pg, err := ConnectToPostgres()
	if err != nil {
		jeagerLogger.Error("Could not connect to database: " + err.Error())
	}

	enc := encoder.NewEncodeClient("encoder", client.DefaultClient)


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
	//encodePublisher := micro.NewPublisher("encoder_pubsub", srv.Client())
	r := rabbitmq.NewBroker(broker.Addrs("amqp://admin:password@rabbit-rabbitmq:5672"))

	if err := r.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := r.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}



	repo := &UploadRepository{s3, pg, tracer, enc, r}

	// Register handler
	pb.RegisterUploadHandler(srv.Server(), &service{repo, tracer, zapLogger})

	// Run the server
	if err := srv.Run(); err != nil {
		jeagerLogger.Error(err.Error())
	}
}