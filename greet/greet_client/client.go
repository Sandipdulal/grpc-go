package main

import (
	"context"
	"fmt"
	"github.com/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
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
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadline(c, 1)
	doUnaryWithDeadline(c, 5)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, duration time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Wick",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration*time.Second)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.DeadlineExceeded {
				fmt.Println("request timeout")
			} else {
				log.Fatalf("unexpected error occured: %v\n", respErr)
			}
		} else {
			log.Fatalf("error calling GreetWithDeadline: %v", err)
		}
		return
	}

	fmt.Printf("response from GreetWithDeadline: %v \n", res.GetResult())

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.GreetEveryoneRequest{
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

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error calling GreetEveryone rpc:%v", err)
	}

	waitChan := make(chan struct{})

	//send multiple requests to server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v \n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error sending stream request: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error closing send channel: %v", err)
		}
	}()

	//receive multiple response from server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				close(waitChan)
				log.Fatalf("error receiving response: %v", err)
			}
			fmt.Printf("Received: %v \n", res.GetResult())
		}
		close(waitChan)
	}()
	<-waitChan

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
