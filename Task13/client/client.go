package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-prog/gen/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	client := pb.NewTestApiClient(conn)

	resp, err := client.Echo(context.Background(), &pb.ResponseRequest{Msg: "Hello World!!!"})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp)
	fmt.Println(resp.Msg)
}