package main

import (
	"context"
	"fmt"
	"github.com/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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
	log.Printf("PrimeNumberDecomposition request from client")
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

func (*server) ComputeAverage(request calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Printf("ComputeAverage request from client")
	sum := int64(0)
	count := 0
	for {
		streamReq, err := request.Recv()
		if err == io.EOF {
			result := float64(sum) / float64(count)
			return request.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error receiving stream request: %v \n", err)
		}
		sum += streamReq.GetNumber()
		count++
		fmt.Println(sum, count)
	}

}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Printf("FindMaximum request from client")
	max := int32(0)
	for {
		streamReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error receiving request: %v", err)
			return err
		}
		num := streamReq.GetNumber()
		if num > max {
			max = num
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: num,
			})
			if err != nil {
				log.Fatalf("error sending response: %v", err)
			}
		}
	}

}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Printf("SquareRoot request from client")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("received a negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
