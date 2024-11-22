# Razorpay Microservice

This microservice provides an interface to interact with the Razorpay API for creating orders, verifying payments, and fetching orders. It uses gRPC for communication and the Razorpay Go SDK for payment-related operations.

## Features

- **CreateOrder**: Create an order on Razorpay with dynamic receipt IDs.
- **VerifyPayment**: Verify the payment signature for an order.
- **FetchOrders**: Fetch orders from Razorpay with options for pagination, timestamp filtering, and receipt-based filtering.

## Requirements

- Go 1.22 or later
- Razorpay API Key & Secret (set via environment variables)
- `.env` file (for managing sensitive keys)

## Setup

### 1. Clone the Repository

```bash

git clone https://github.com/yourusername/razorpay-microservice.git
cd razorpay-microservice

```
### 2. Install Dependencies

```bash

go mod tidy

```
### 3. Setup Environment Variables

```bash

RAZORPAY_KEY=your_razorpay_api_key
RAZORPAY_SECRET=your_razorpay_api_secret

```
### 4. Run the Microservice

```bash
go run cmd/main.go

```
