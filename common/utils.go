package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"os"
)

// HMAC calculation function
func CalculateHMAC(orderID, paymentID string) []byte {
	apiSecret := os.Getenv("RAZORPAY_SECRET")
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(orderID + "|" + paymentID))
	return h.Sum(nil)
}
