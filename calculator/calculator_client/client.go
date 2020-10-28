package main

import (
	"context"
	"github.com/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect to grpc server: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
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
