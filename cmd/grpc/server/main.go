package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gcarrenho/guidemysteps/internal/routing"
	"github.com/gcarrenho/guidemysteps/internal/routing/routinggrpc"
	"github.com/gcarrenho/guidemysteps/internal/translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	routingRepo := routing.NewRoutingRepo("https://routing.openstreetmap.de", translateSvc)
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
