package main

import (
	"log"
	"net"
	"os"
	"time"

	"local/github_info_system/collector/internal/delivery/grpc"
	"local/github_info_system/collector/internal/repository/github"
	"local/github_info_system/collector/internal/usecase"
	"local/github_info_system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	githubClient := github.NewClient(10 * time.Second)
	repoUseCase := usecase.NewRepositoryUseCase(githubClient)

	grpcServer := grpc.NewServer()
	collectorServer := grpc.NewServer(repoUseCase)

	proto.RegisterCollectorServiceServer(grpcServer, collectorServer)
	// Enable reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("Collector service starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
