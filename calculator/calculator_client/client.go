package main

import (
	"context"
	"fmt"
	"github.com/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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
	//doComputeAverage(c)
	doFindMaximum(c)

}

func doFindMaximum(c calculatorpb.CalculatorServiceClient) {
	requests := []int32{10, 12, 30, 14, 22, 50, 3, 6, 79}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error calling FindMaximum rpc: %v", err)
	}

	closeChan := make(chan struct{})

	//send multiple request
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v \n", req)
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: req,
			})
			if err != nil {
				log.Fatalf("error sending rpc request: %v \n", err)
			}
			time.Sleep(1 * time.Second)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("error closing send stream: %v \n", err)
		}
	}()

	//receive multiple response
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				close(closeChan)
				log.Fatalf("error receiving message: %v \n", err)
			}
			maximum := res.GetResult()
			fmt.Printf("Received: %v \n", maximum)
		}
		close(closeChan)
	}()

	<-closeChan
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
