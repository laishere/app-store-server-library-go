package appstore

import (
	"testing"
)

// Test promotional offer V2 signature creation
func TestPromotionalOfferV2SignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", stringPtr("999"))
	assertNoError(t, err, "Failed to create signature")

	// Verify the JWT was created
	assertTrue(t, signature != "", "Signature is not empty")

	// Parse and verify the JWT structure
	header, payload, err := decodeJWTWithoutVerification(signature)
	assertNoError(t, err, "Failed to decode JWT")

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
	assertTrue(t, payload["nonce"] != nil && payload["nonce"] != "", "Nonce is present")

	// Verify iat exists
	assertNotNil(t, payload["iat"], "IAT")
}

// Test promotional offer V2 without transaction ID
func TestPromotionalOfferV2SignatureCreator_NoTransactionId(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", nil)
	assertNoError(t, err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assertNoError(t, err, "Failed to decode JWT")

	// Transaction ID should not be in payload
	assertNil(t, payload["transactionId"], "Transaction ID")
}

// Test promotional offer V2 with missing required fields
func TestPromotionalOfferV2SignatureCreator_MissingProductId(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	_, err = creator.CreateSignature("", "OFFER123", nil)
	assertError(t, err, "Expected error for empty product ID")
}

func TestPromotionalOfferV2SignatureCreator_MissingOfferIdentifier(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	_, err = creator.CreateSignature("com.example.product", "", nil)
	assertError(t, err, "Expected error for empty offer identifier")
}
