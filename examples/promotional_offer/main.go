package main

import (
	"fmt"
	"os"

	appstore "github.com/laishere/app-store-server-library-go"
)

func main() {
	// Loading your private key
	signingKey, err := os.ReadFile("SubscriptionKey_ABCD123456.p8")
	if err != nil {
		fmt.Printf("Failed to read key: %v\n", err)
		return
	}

	keyId := "ABCD123456"
	issuerId := "57246542-96fe-1a63-e053-0824d011072a"
	bundleId := "com.example.app"

	// Create the signature creator for Promotional Offer V2
	creator, err := appstore.NewPromotionalOfferV2SignatureCreator(
		signingKey,
		keyId,
		issuerId,
		bundleId,
	)
	if err != nil {
		panic(err)
	}

	productID := "my_subscription_product"
	offerIdentifier := "my_offer_id"

	// transactionID is optional for V2 signatures
	// If you have an existing transaction you want to link the offer to, provide it here.
	var transactionID *string = nil

	signature, err := creator.CreateSignature(productID, offerIdentifier, transactionID)
	if err != nil {
		fmt.Printf("Signature creation failed: %v\n", err)
		return
	}

	fmt.Printf("Generated JWS Signature: %s\n", signature)
}
