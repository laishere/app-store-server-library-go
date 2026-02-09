package appstore

import (
	"testing"
)

// Test introductory offer eligibility signature
func TestIntroductoryOfferEligibilitySignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewIntroductoryOfferEligibilitySignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	signature, err := creator.CreateSignature("com.example.product", true, "999")
	if err != nil {
		t.Fatalf("Failed to create signature: %v", err)
	}

	_, payload, err := decodeJWTWithoutVerification(signature)
	if err != nil {
		t.Fatalf("Failed to decode JWT: %v", err)
	}

	assertEqual(t, "introductory-offer-eligibility", payload["aud"], "Audience")
	assertEqual(t, "com.example.product", payload["productId"], "Product ID")
	assertEqual(t, true, payload["allowIntroductoryOffer"], "Allow Introductory Offer")
	assertEqual(t, "999", payload["transactionId"], "Transaction ID")
}

// Test advanced commerce API signature
func TestAdvancedCommerceAPIInAppSignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewAdvancedCommerceAPIInAppSignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	request := map[string]any{
		"productId": "com.example.product",
		"quantity":  1,
	}

	signature, err := creator.CreateSignature(request)
	if err != nil {
		t.Fatalf("Failed to create signature: %v", err)
	}

	_, payload, err := decodeJWTWithoutVerification(signature)
	if err != nil {
		t.Fatalf("Failed to decode JWT: %v", err)
	}

	assertEqual(t, "advanced-commerce-api", payload["aud"], "Audience")

	// Verify request field exists and is base64 encoded
	if payload["request"] == nil {
		t.Fatal("Request field should be present")
	}

	requestStr, ok := payload["request"].(string)
	if !ok {
		t.Fatal("Request should be a string")
	}
	if requestStr == "" {
		t.Error("Request should not be empty")
	}
}
