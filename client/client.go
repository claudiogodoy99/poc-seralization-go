package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "poc-serialization/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr   = flag.String("addr", "localhost:50001", "the address to connect to")
	action = flag.String("action", "inc", "")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSerializedServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if *action == "cc" {
		r, err := c.ParentCC(ctx, &pb.CC{Counter: 0})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %t", r.Ok)
	} else if *action == "nocc" {
		r, err := c.UnaryCall(ctx, &pb.UnaryRequest{ValueToIncrement: 1000})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("Greeting: %t", r.Ok)
	} else if *action == "inc" {
		r, err := c.ParentUnaryCall(ctx, &pb.UnaryRequest{ValueToIncrement: 4})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("Greeting: %t", r.Ok)
	}

}
