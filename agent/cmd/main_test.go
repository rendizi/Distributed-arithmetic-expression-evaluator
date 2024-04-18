package main_test

import (
	"context"
	"google.golang.org/grpc"
	"testing"

	pb "github.com/rendizi/Distributed-arithmetic-expression-evaluator/proto"
)

var ports = []string{"localhost:5000", "localhost:5001", "localhost:5002"}

func TestAv(t *testing.T) {
	for _, port := range ports {
		conn, err := grpc.Dial(port, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewAgentServiceClient(conn)

		request := &pb.AvRequest{}

		response, err := client.Av(context.Background(), request)
		if err != nil {
			t.Fatalf("Av RPC failed: %v", err)
		}

		if !response.Result {
			t.Error("Agent not available")
		}
	}
}

func TestOp(t *testing.T) {
	for _, port := range ports {
		conn, err := grpc.Dial(port, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewAgentServiceClient(conn)

		request := &pb.OpRequest{
			A:        10,
			B:        5,
			Operator: "+",
			Time:     1,
		}

		response, err := client.Op(context.Background(), request)
		if err != nil {
			t.Fatalf("Op RPC failed: %v", err)
		}

		expected := float32(15)
		if response.Result != expected {
			t.Errorf("Expected %f, got %f", expected, response.Result)
		}
	}
}
