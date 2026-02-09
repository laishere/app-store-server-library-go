package appstore

import (
	"testing"
)

// Test promotional offer V2 signature creation
func TestPromotionalOfferV2SignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", stringPtr("999"))
	if err != nil {
		t.Fatalf("Failed to create signature: %v", err)
	}

	// Verify the JWT was created
	if signature == "" {
		t.Fatal("Signature should not be empty")
	}

	// Parse and verify the JWT structure
	header, payload, err := decodeJWTWithoutVerification(signature)
	if err != nil {
		t.Fatalf("Failed to decode JWT: %v", err)
	}

	// Verify header
	assertEqual(t, "ES256", header["alg"], "Algorithm")
	assertEqual(t, TEST_KEY_ID, header["kid"], "Key ID")

	// Verify payload
	assertEqual(t, TEST_BUNDLE_ID, payload["bid"], "Bundle ID")
	assertEqual(t, TEST_ISSUER_ID, payload["iss"], "Issuer ID")
	assertEqual(t, "promotional-offer", payload["aud"], "Audience")
	assertEqual(t, "com.example.product", payload["productId"], "Product ID")
	assertEqual(t, "OFFER123", payload["offerIdentifier"], "Offer Identifier")
	assertEqual(t, "999", payload["transactionId"], "Transaction ID")

	// Verify nonce exists
	if payload["nonce"] == nil || payload["nonce"] == "" {
		t.Error("Nonce should be present")
	}

	// Verify iat exists
	if payload["iat"] == nil {
		t.Error("IAT should be present")
	}
}

// Test promotional offer V2 without transaction ID
func TestPromotionalOfferV2SignatureCreator_NoTransactionId(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", nil)
	if err != nil {
		t.Fatalf("Failed to create signature: %v", err)
	}

	_, payload, err := decodeJWTWithoutVerification(signature)
	if err != nil {
		t.Fatalf("Failed to decode JWT: %v", err)
	}

	// Transaction ID should not be in payload
	if payload["transactionId"] != nil {
		t.Error("Transaction ID should not be present when nil")
	}
}

// Test promotional offer V2 with missing required fields
func TestPromotionalOfferV2SignatureCreator_MissingProductId(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	_, err = creator.CreateSignature("", "OFFER123", nil)
	if err == nil {
		t.Fatal("Expected error for empty product ID")
	}
}

func TestPromotionalOfferV2SignatureCreator_MissingOfferIdentifier(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	if err != nil {
		t.Fatalf("Failed to read signing key: %v", err)
	}

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	if err != nil {
		t.Fatalf("Failed to create signature creator: %v", err)
	}

	_, err = creator.CreateSignature("com.example.product", "", nil)
	if err == nil {
		t.Fatal("Expected error for empty offer identifier")
	}
}
