package main

import (
	"context"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
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
