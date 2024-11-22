package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"razorpay-microservice/common"
	"razorpay-microservice/pb"
	"time"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
}

func GetDefaultCurrency(currency *string) string {
	if currency != nil && *currency != "" {
		return *currency
	}
	return "INR"
}

func (s *PaymentServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	client, err := common.RazoryClient()
	if req.Amount <= 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"Amount must be greater than zero",
		)
	}
	receiptID := uuid.New().String()

	orderData := map[string]interface{}{
		"amount":   req.Amount,
		"currency": GetDefaultCurrency(req.Currency),
		"receipt":  receiptID, // Use dynamic receipt ID
	}
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		// Handle errors returned from Razorpay (e.g., network issues, API errors)
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Failed to create order with Razorpay: %v", err),
		)
	}
	return &pb.CreateOrderResponse{
		OrderId:   order["id"].(string),
		ReceiptId: receiptID,
		Amount:    int32(order["amount"].(float64)),
		Currency:  order["currency"].(string),
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *PaymentServiceServer) VerifyPayment(ctx context.Context, req *pb.VerifyPaymentRequest) (*pb.VerifyPaymentResponse, error) {
	if req.RazorpayOrderId == "" || req.RazorpayPaymentId == "" || req.RazorpaySignature == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"Missing required fields: RazorpayOrderId, RazorpayPaymentId, or RazorpaySignature",
		)
	}

	expectedSignature := fmt.Sprintf("%x", common.CalculateHMAC(req.RazorpayOrderId, req.RazorpayPaymentId))

	// Verify the signature
	if expectedSignature != req.RazorpaySignature {
		return nil, status.Error(
			codes.InvalidArgument,
			"Invalid payment signature. Verification failed.",
		)
	}

	return &pb.VerifyPaymentResponse{
		Valid:   true,
		Message: "Payment verified successfully",
	}, nil
}

func (s *PaymentServiceServer) FetchOrders(ctx context.Context, req *pb.FetchOrdersRequest) (*pb.FetchOrdersResponse, error) {
	// Check if the count is missing or invalid
	if req.Count == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"Count is required and must be greater than zero",
		)
	}
	// Prepare options based on the incoming request
	options := map[string]interface{}{
		"count": req.Count, // Number of orders to fetch
		"skip":  req.Skip,  // Number of orders to skip
	}

	if req.From != nil && *req.From > 0 {
		options["from"] = *req.From // Timestamp after which orders were created
	}
	if req.To != nil && *req.To > 0 {
		options["to"] = *req.To // Timestamp before which orders were created
	}

	// Apply receipt filter if provided
	if req.Receipt != nil && *req.Receipt != "" {
		options["receipt"] = (*req.Receipt)
	}

	client, err := common.RazoryClient()
	// Fetch the orders from Razorpay API

	orders, err := client.Order.All(options, nil)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Failed to fetch orders with Razorpay: %v", err),
		)
	}

	// Convert orders into response format
	var orderList []*pb.Order
	for _, order := range orders["items"].([]interface{}) {
		orderData := order.(map[string]interface{})
		orderList = append(orderList, &pb.Order{
			OrderId:   orderData["id"].(string),
			Amount:    int32(orderData["amount"].(float64)),
			Currency:  orderData["currency"].(string),
			ReceiptId: orderData["receipt"].(string),
			Status:    orderData["status"].(string),
			CreatedAt: int32(orderData["created_at"].(float64)),
		})
	}

	// Return the response with a list of orders
	return &pb.FetchOrdersResponse{
		Orders: orderList,
	}, nil
}
