package main

import (
	"context"
	"fmt"
	"io"

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
	// doUnary(c)
	// doStreaming(c)
	doClientStreaming(c)
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

func doStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Created client: %f", c)

	req := &calculatorpb.PrimeRequest{
		Number: 1203212,
	}
	stream, err := c.FindPrime(context.Background(), req)

	if err != nil {
		log.Fatalf("prime error: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}
		log.Print(msg.GetResult())
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.AvgRequest{
		{Num: 1},
		{Num: 2},
		{Num: 3},
		{Num: 4},
		{Num: 5},
		{Num: 6},
		{Num: 7},
		{Num: 8},
		{Num: 9},
		{Num: 10},
	}

	stream, err := c.FindAvg(context.Background())
	if err != nil {
		log.Fatalf("error calling FindAvg : %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receving response from FindAvg : %v", err)
	}
	log.Printf("Find Avg Response : %v", res.GetResult())

}
