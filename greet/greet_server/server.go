package main

import (
	"context"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function invoked by client: %v \n", request)
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()
	res := &greetpb.GreetResponse{
		Result: "Hello! " + firstName + " " + lastName,
	}
	return res, nil
}

func (s *Server) GreetManyTimes(request *greetpb.GreetMayTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes stream invoked by client: %v \n", request)
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello! " + firstName + " " + lastName + ":" + strconv.Itoa(i),
		}
		err := stream.Send(res)
		if err != nil {
			log.Fatalf("unable to send response stream: %v \n", err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to start tcp listner: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
