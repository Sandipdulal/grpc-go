package main

import (
	"context"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
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

func (s *Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("LongGreet call invoked by client")
	result := ""
	for {
		streamReq, err := stream.Recv()
		res := &greetpb.LongGreetResponse{
			Result: result,
		}
		if err == io.EOF {
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("error reading stream request: %v", err)
		}
		result += "Hello! " + streamReq.GetGreeting().GetFirstName() + " " + streamReq.GetGreeting().GetLastName() + "."
	}
}

func (s *Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone call invoked by client")
	for {
		streamReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error receiving stream request: %v", err)
			return err
		}
		result := "Hello, " + streamReq.GetGreeting().GetFirstName() + " " + streamReq.GetGreeting().GetLastName() + "!"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("error sending greet everyone response: %v", err)
			return err
		}
	}
}

func (s *Server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Println("GreetWithDeadline call invoked by client")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "request cancelled by client")
		}
		time.Sleep(1 * time.Second)
	}
	result := "Hello," + req.GetGreeting().GetFirstName() + " " + req.GetGreeting().GetLastName() + "!"
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to start tcp listner: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if err != nil {
		log.Fatalf("error loading server certificates: %v \n", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
