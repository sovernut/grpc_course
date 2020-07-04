package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/labstack/gommon/log"
	"github.com/sovernut/grpc_course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalcRequest) (*calculatorpb.CalcResponse, error) {
	fmt.Printf("Sum function was invoked with %v", req)
	first := req.GetSum().GetFirst()
	second := req.GetSum().GetSecond()

	result := first + second
	res := &calculatorpb.CalcResponse{
		Result: result,
	}
	return res, nil
}

func (*server) FindPrime(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_FindPrimeServer) error {
	fmt.Printf("FindPrime function was invoked with %v", req)
	number := req.GetNumber()
	var k int64 = 2
	N := number
	for N > 1 {
		if N%k == 0 { // if k evenly divides into N
			res := &calculatorpb.CalcResponse{
				Result: int64(k),
			}
			stream.Send(res)
			N = N / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) FindAvg(stream calculatorpb.CalculatorService_FindAvgServer) error {
	fmt.Printf("FindAvg function was invoked with %v", stream)

	var n int64
	var sum int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			ans := float32(sum) / float32(n)
			return stream.SendAndClose(&calculatorpb.AvgResponse{
				Result: ans,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream :%v", err)
		}

		recvNum := req.GetNum()
		sum += recvNum
		n++
	}
	return nil
}
func main() {
	fmt.Println("hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
