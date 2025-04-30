package main

import (
	"context"
	"log"
	"time"

	"github.com/xavicci/gRPC-Student-Service/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := testpb.NewTestServiceClient(cc)
	DoUnary(c)
	DoClientStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest: %v", err)
	}
	log.Printf("response from GetTest: %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q1",
			Question: "What is the capital of France?",
			Answer:   "Paris",
			TestId:   "t1",
		},
		{
			Id:       "q2",
			Question: "What is the capital of Germany?",
			Answer:   "Berlin",
			TestId:   "t1",
		},
		{
			Id:       "q3",
			Question: "What is the capital of Italy?",
			Answer:   "Rome",
			TestId:   "t1",
		},
	}
	stream, err := c.SetQuestion(context.Background())
	if err != nil {
		log.Fatalf("error while calling SetQuestion: %v", err)
	}

	for _, question := range questions {
		log.Printf("sending question: %v", question)
		stream.Send(question)
		time.Sleep(1 * time.Second)
	}
	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	log.Printf("response from SetQuestion: %v", msg)
}
