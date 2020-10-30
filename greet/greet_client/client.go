package main

import (
	"context"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Wick",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Cena",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Van",
				LastName:  "Dam",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling LongGreet: %v", err)
	}
	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			log.Fatalf("error sending stream request: %v", err)
		}
	}
	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receving stream response: %v", err)
	}
	log.Printf("LongGreet response: %v \n", result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetMayTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Wick",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling GreetManyTimes: %v \n", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			//end of stream
			log.Printf("stream ended")
			break
		}
		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}
		log.Printf("stream response:%v \n", res.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Wick",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to invoke Greet rpc : %v \n", err)
	}

	log.Printf("response from Greet rpc :%v \n", res)
}
