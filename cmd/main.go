package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"razorpay-microservice/common"
	"razorpay-microservice/pb"
	"razorpay-microservice/pkg/payment"
)

func main() {
	// Load environment variables
	common.LoadEnv()
	// Listen on port 50051
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(common.AuthUnaryInterceptor),
	)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Register the PaymentServiceServer with the gRPC server
	pb.RegisterPaymentServiceServer(grpcServer, &payment.PaymentServiceServer{})

	// Log when the server starts
	log.Println("gRPC server starting on port :50051...")

	// Start the gRPC server to listen for incoming requests
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
