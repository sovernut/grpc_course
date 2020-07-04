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
	// doClientStreaming(c)
	doBiStreaming(c)
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

func doBiStreaming(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.MaxRequest{
		{Num: 1},
		{Num: 10},
		{Num: 15},
		{Num: 4},
		{Num: 5},
		{Num: 20},
		{Num: 7},
		{Num: 8},
		{Num: 9},
		{Num: 10},
	}

	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("error calling findmax : %v", err)
	}
	waitc := make(chan int)

	// use goroutine to send each request
	go func() {
		for _, req := range requests {
			if err := stream.Send(req); err != nil {
				log.Fatalf("error sending max: %v", err)
			}
		}
		// call when send msg completed
		stream.CloseSend()
	}()

	go func() {
		// loop forever
		for {
			// recieve each msg from server
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error recieving max: %v", err)
				break
			}
			// get result from msg
			log.Printf("Receiving max: %v", msg.GetMax())
		}
		// close channel waitc
		close(waitc)
	}()

	// block until it has value or close.
	<-waitc

}
