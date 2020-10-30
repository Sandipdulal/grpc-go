package main

import (
	"context"
	"github.com/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect to grpc server: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doSum(c)
	//doPrimeDecomposition(c)
	doComputeAverage(c)

}

func doComputeAverage(c calculatorpb.CalculatorServiceClient) {
	requests := []int64{1, 2, 3, 4}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error calling ComputeAverage rpc: %v \n", err)
	}

	for _, req := range requests {
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: req,
		})
		if err != nil {
			log.Fatalf("error sending stream request: %v \n", err)
		}
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response from server: %v \n", err)
	}
	log.Printf("ComputeAverage response:%v \n", result.GetResult())

}

func doPrimeDecomposition(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	streamResponse, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling PrimeNumberDecomposition: %v \n", err)
	}
	for {
		res, err := streamResponse.Recv()
		if err == io.EOF {
			log.Printf("end of stream")
			break
		}
		if err != nil {
			log.Fatalf("error receiving stream response: %v \n", err)
		}
		log.Printf("Stream result: %v \n", res.GetResult())
	}
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to invoke Sum rpc: %v \n", err)
	}
	log.Printf("Sum response: %v \n", res)
}
