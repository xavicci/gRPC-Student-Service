package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/xavicci/gRPC-Student-Service/models"
	"github.com/xavicci/gRPC-Student-Service/repository"
	"github.com/xavicci/gRPC-Student-Service/studentpb"
	"github.com/xavicci/gRPC-Student-Service/testpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Test not found: %v", err)
	}
	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}
	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to set test: %v", err)
	}
	return &testpb.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetQuestion(stream testpb.TestService_SetQuestionServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}
		question := &models.Question{
			Id:       msg.Id,
			Answer:   msg.Answer,
			Question: msg.Question,
			TestId:   msg.TestId,
		}
		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}
		enrollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.TestId,
		}
		err = s.repo.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.TestId)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to get students per test: %v", err)
	}
	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  int32(student.Age),
		}
		err := stream.Send(student)
		time.Sleep(time.Second * 2)
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to send student: %v", err)
		}
	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "t1")
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to get questions per test: %v", err)
	}
	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}
		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				return status.Errorf(codes.Internal, "Failed to send question: %v", err)
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to receive answer: %v", err)
		}
		log.Println("Answer: ", answer.Answer)
	}
}
