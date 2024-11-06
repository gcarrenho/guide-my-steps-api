package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gcarrenho/guidemysteps/internal/routing"
	"github.com/gcarrenho/guidemysteps/internal/routing/routinggrpc"
	"github.com/gcarrenho/guidemysteps/internal/translator"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "path/to/server.crt", "The file containing the CA root certificate")
	keyFile    = flag.String("key_file", "path/to/server.key", "The file containing the server's private key")
	port       = flag.Int("port", 50051, "The server port")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	bundle     = i18n.NewBundle(language.English)
)

func main() {
	flag.Parse()
	translateRepo := translator.NewI18nRepo(bundle)
	translateSvc := translator.NewTranslationService(translateRepo)

	routingRepo := routing.NewOpenStreetMapProvider("https://routing.openstreetmap.de", translateSvc)
	routingSvc := routing.NewRoutingComponentImpl(routingRepo)
	/*if err != nil {
		log.Fatalf("Failed to create feature repository: %v", err)
	}*/
	//featureSvc := routing.NewRoutingComponentImpl(repo)

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("failed to generate credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	logger, _ := zap.NewProduction() // Puedes usar zap.NewDevelopment() para ambiente de desarrollo
	defer logger.Sync()

	//logger := zap.L().Named("server") // create a log
	requestID := uuid.New().String()
	requestLogger := logger.With(zap.String("request_id", requestID))

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64(
					"grpc.time_ns",
					duration.Nanoseconds(),
				)
			},
		),
	} // Config of options of logger. In thiss case the duration in nanoseconds

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // Config of th trace to always show the request
	err := view.Register(ocgrpc.DefaultServerViews...)                    // View register of opencensus
	if err != nil {
		//logger.Fatal().Err(err).Msg("Failed to register views")
	}

	opts = append(opts,
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_zap.StreamServerInterceptor(requestLogger, zapOpts...),
				grpc_auth.StreamServerInterceptor(authenticate),
			)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(requestLogger, zapOpts...),
			grpc_auth.UnaryServerInterceptor(authenticate),
		)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)

	gRPCServer, err := routinggrpc.NewGRPCServer(&routinggrpc.Config{FeatureSvc: routingSvc}, opts...)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Println("Listening and serving HTTP on :50051")
	grpcLn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := gRPCServer.Serve(grpcLn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.New(
			codes.Unknown,
			"couldn't find peer info",
		).Err()
	}

	if peer.AuthInfo == nil {
		return context.WithValue(ctx, subjectContextKey{}, ""), nil
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	ctx = context.WithValue(ctx, subjectContextKey{}, subject)

	return ctx, nil
}

type subjectContextKey struct{}
