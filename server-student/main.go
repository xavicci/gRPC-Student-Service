package main

import (
	"log"
	"net"

	"github.com/xavicci/gRPC-Student-Service/database"
	"github.com/xavicci/gRPC-Student-Service/server"
	"github.com/xavicci/gRPC-Student-Service/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/mibase?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	server := server.NewStudentServer(repo)

	grpcServer := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
