package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/http"
	grpchandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port: %d", port)
	registry, err := consul.NewRegistry("localhost:8500")

	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)

	ctx := context.Background()

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {

		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}

	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratinggateway := ratinggateway.New(registry)

	ctrl := movie.New(ratinggateway, metadataGateway)

	// gRPC protocol:
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}

	// HTTP protocol:
	// h := httphandler.New(ctrl)

	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	// if err := http.ListenAndServe(":8083", nil); err != nil {
	// 	panic(err)
	// }

}
