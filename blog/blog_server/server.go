package main

import (
	"context"
	"fmt"
	"github.com/grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

var collection *mongo.Collection

func main() {
	//set logging flags to print lines in case the program crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog server started...")

	fmt.Println("connecting to mongodb server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("unable to connect to mongodb server")
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to start tcp listener: %v", err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	//shutdown gracefully , wait for ctrl C signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//block until signal is received
	<-ch
	fmt.Println("stopping the server...")
	s.Stop()
	fmt.Println("stopping the listener")
	lis.Close()
	fmt.Println("closing mongodb client")
	client.Disconnect(context.TODO())
	fmt.Println("end of program!")
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("create blog request")
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot convert to OID"))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("read blog request")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot parse blog id: %v", err))
	}

	data := &blogItem{}

	filter := bson.M{"_id": oid}
	res := collection.FindOne(context.Background(), filter)
	decodeErr := res.Decode(data)
	if decodeErr != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog with id specified: %v", decodeErr))
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot parse blog id: %v", err))
	}

	data := &blogItem{}

	filter := bson.M{"_id": oid}
	res := collection.FindOne(context.Background(), filter)
	decodeErr := res.Decode(data)
	if decodeErr != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog with id specified: %v", decodeErr))
	}

	data.Title = blog.GetTitle()
	data.Content = blog.GetContent()
	data.AuthorID = blog.GetAuthorId()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("unable to update item in mongodb: %v", updateErr))
	}

	fmt.Println("successfully updated record with updateid")

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot parse blog id: %v", err))
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot delete blog with id specified: %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog with id specified: %v", blogID))

	}
	fmt.Printf("successfully deleted record with id:%s \n", blogID)
	return &blogpb.DeleteBlogResponse{BlogId: blogID}, nil
}

func (s *server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("error calling Find on collection")
		return status.Errorf(codes.Internal, fmt.Sprintf("unknown internal error: %v", err))
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &blogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("cannot decode cursor data: %v", err))
		}
		sendErr := stream.Send(&blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorID,
				Title:    data.Title,
				Content:  data.Content,
			},
		})

		if sendErr != nil {
			log.Fatalf("error sending response: %v", sendErr)
			return err
		}

		if cursor.Err() != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("unknown internal error: %v", err))
		}

	}
	return nil
}
