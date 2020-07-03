package main

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/sovernut/grpc_course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Created client: %f", c)

	req := &calculatorpb.CalcRequest{
		Sum: &calculatorpb.Sum{
			First:  7,
			Second: 3,
		},
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Greet response: %v", err)
	}

	log.Printf("Response from Calculator %v", res.Result)
}
