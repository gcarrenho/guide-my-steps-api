package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	api "github.com/gcarrenho/guidemysteps/api/v1/mysteps"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/data"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

// printFeature gets the feature for the given point.
func printFeature(client api.MyStepsClient, request *api.MyStepsRequest) {
	//log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetRoute(ctx, request)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}
	log.Println("respuesta ", feature)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := api.NewMyStepsClient(conn)

	// Looking for a valid feature
	requets := &api.MyStepsRequest{
		Start: &api.LatLng{
			Longitude: 2.2630370544817957,
			Latitude:  41.93002134481184,
		},
		Destination: &api.LatLng{
			Longitude: 2.256949207974631,
			Latitude:  41.93002134481184,
		},
	}
	start := time.Now()
	printFeature(client, requets)
	end := time.Now().Sub(start)
	fmt.Println("time in milisecond ", end.Milliseconds())

}
