package main

import (
	"context"
	"fmt"
	"github.com/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Blog client started...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	doCreateBlog(c)
	//doReadBlog(c)
	//doUpdateBlog(c)
	//doDeleteBlog(c)
	doListBlog(c)
}

func doListBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.ListBlogRequest{}
	res, err := c.ListBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling ListBlog rpc: %v", err)
	}

	for {
		stream, err := res.Recv()
		if err == io.EOF {
			fmt.Println("end of stream")
			break
		}
		if err != nil {
			fmt.Printf("error receiving stream response: %v \n", err)
		}
		fmt.Printf("ListBlog response:%v \n", stream.GetBlog())
	}

}

func doDeleteBlog(c blogpb.BlogServiceClient) {

	req := &blogpb.DeleteBlogRequest{
		BlogId: "5fa879746320e05dbf8b92e5",
	}
	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling DeleteBlog rpc: %v", err)
	}
	fmt.Printf("DeleteBlog response:%v \n", res.GetBlogId())
}

func doUpdateBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5fa879746320e05dbf8b92e5",
			AuthorId: "John12345",
			Title:    "Updated record",
			Content:  "This record has been updated",
		},
	}

	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling UpdateBlog rpc: %v", err)
	}
	fmt.Printf("UpdateBlog response:%v \n", res.GetBlog())
}

func doReadBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.ReadBlogRequest{
		BlogId: "5fa879746310e05dbf8b92e5",
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling ReadBlog rpc: %v", err)
	}
	fmt.Printf("ReadBlog response:%v \n", res.GetBlog())

}

func doCreateBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "John123",
			Title:    "TestBlog",
			Content:  "This is the test rpc blog request",
		},
	}
	blogRes, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling CreateBlog rpc: %v", err)
	}

	fmt.Printf("Blog response:%v \n", blogRes.GetBlog())
}
