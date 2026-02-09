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

	creator, err := appstore.NewIntroductoryOfferEligibilitySignatureCreator(
		signingKey,
		keyId,
		issuerId,
		bundleId,
	)
	if err != nil {
		panic(err)
	}

	productID := "subscription_product"
	allowIntroductoryOffer := true
	transactionID := "1000000001" // An existing transaction ID for the user

	signature, err := creator.CreateSignature(productID, allowIntroductoryOffer, transactionID)
	if err != nil {
		fmt.Printf("Signature creation failed: %v\n", err)
		return
	}

	fmt.Printf("Generated Introductory Offer Eligibility Signature: %s\n", signature)
}
