package main

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/sovernut/grpc_course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Created client: %f", c)

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Greet response: %v", err)
	}

	log.Printf("Response from Greet %v", res.Result)
}
