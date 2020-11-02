package main

import (
	"context"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
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
