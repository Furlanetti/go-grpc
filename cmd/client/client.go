package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/furlanetti/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
  AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "test@test.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make grpc request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "test@test.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make grpc request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status:", stream.Status, "- ", stream.User)
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "f1",
			Name:  "Felipe1",
			Email: "f1@felipe.com",
		},
		&pb.User{
			Id:    "f2",
			Name:  "Felipe2",
			Email: "f2@felipe.com",
		},
		&pb.User{
			Id:    "f3",
			Name:  "Felipe3",
			Email: "f3@felipe.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating requests: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating requests: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "f1",
			Name:  "Felipe1",
			Email: "f1@felipe.com",
		},
		&pb.User{
			Id:    "f2",
			Name:  "Felipe2",
			Email: "f2@felipe.com",
		},
		&pb.User{
			Id:    "f3",
			Name:  "Felipe3",
			Email: "f3@felipe.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Println("Recebendo user %v com status %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
