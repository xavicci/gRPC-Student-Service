package server

import (
	"context"

	"github.com/xavicci/gRPC-Student-Service/models"
	"github.com/xavicci/gRPC-Student-Service/repository"
	"github.com/xavicci/gRPC-Student-Service/studentpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	student, err := s.repo.GetStudent(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Student not found: %v", err)
	}
	return &studentpb.Student{
		Id:   &student.Id,
		Name: &student.Name,
		Age:  &student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {

	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}
	err := s.repo.SetStudent(ctx, student)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to set student: %v", err)
	}
	return &studentpb.SetStudentResponse{
		Id: &student.Id,
	}, nil

}
