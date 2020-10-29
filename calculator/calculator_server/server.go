package main

import (
	"context"
	"github.com/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (res *calculatorpb.SumResponse, err error) {
	log.Printf("Sum request from client: %v", req)
	sum := req.GetFirstNumber() + req.GetSecondNumber()
	result := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return result, nil

}

func (*server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	k := int64(2)
	n := request.GetNumber()

	for n > 1 {
		if n%k == 0 {
			streamResponse := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			err := stream.Send(streamResponse)
			if err != nil {
				log.Fatalf("error sending stream response: %v \n", err)
			}
			n = n / k
		} else {
			k++
		}
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to listen on tcp port: %v \n", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("unable to start server: %v \n", err)
	}

}
