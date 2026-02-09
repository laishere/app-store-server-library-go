package main

import (
	"fmt"
	"os"

	appstore "github.com/laishere/app-store-server-library-go"
)

func main() {
	signingKey, _ := os.ReadFile("SubscriptionKey_ABCD123456.p8")
	keyId := "ABCD123456"
	issuerId := "57246542-96fe-1a63-e053-0824d011072a"
	bundleId := "com.example.app"

	creator, err := appstore.NewAdvancedCommerceAPIInAppSignatureCreator(
		signingKey,
		keyId,
		issuerId,
		bundleId,
	)
	if err != nil {
		panic(err)
	}

	// This can be any request object required by the Advanced Commerce API
	requestBody := map[string]any{
		"some_field": "some_value",
		"timestamp":  1700000000000,
	}

	signature, err := creator.CreateSignature(requestBody)
	if err != nil {
		fmt.Printf("Signature creation failed: %v\n", err)
		return
	}

	fmt.Printf("Generated Advanced Commerce API Signature: %s\n", signature)
}
