package main

import (
	"log"
	"net/http"
	"os"

	"task2/github_info_system/api_gateway/internal/delivery/grpc"

	httpHandler "task2/github_info_system/api_gateway/internal/delivery/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "collector:50051"
	}

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to collector service
	collectorClient, err := grpc.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to connect to collector: %v", err)
	}
	defer collectorClient.Close()

	// Setup handler
	handler := httpHandler.NewHandler(collectorClient)

	// Setup router
	r := mux.NewRouter()
	
	// API routes
	r.HandleFunc("/api/v1/repos/{owner}/{repo}", handler.GetRepositoryInfo).Methods("GET")
	
	// Swagger routes
	r.HandleFunc("/swagger.json", handler.SwaggerJSON).Methods("GET")
	r.HandleFunc("/swagger", handler.SwaggerUI).Methods("GET")
	r.HandleFunc("/swagger/", handler.SwaggerUI).Methods("GET")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handlerWithCORS := c.Handler(r)

	log.Printf("API Gateway starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handlerWithCORS); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
