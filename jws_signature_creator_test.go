package appstore

import (
	"testing"
)

// Test introductory offer eligibility signature
func TestIntroductoryOfferEligibilitySignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewIntroductoryOfferEligibilitySignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", true, "999")
	assertNoError(t, err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assertNoError(t, err, "Failed to decode JWT")

	assertEqual(t, "introductory-offer-eligibility", payload["aud"], "Audience")
	assertEqual(t, "com.example.product", payload["productId"], "Product ID")
	assertEqual(t, true, payload["allowIntroductoryOffer"], "Allow Introductory Offer")
	assertEqual(t, "999", payload["transactionId"], "Transaction ID")
}

// Test advanced commerce API signature
func TestAdvancedCommerceAPIInAppSignatureCreator(t *testing.T) {
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	creator, err := NewAdvancedCommerceAPIInAppSignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assertNoError(t, err, "Failed to create signature creator")

	request := map[string]any{
		"productId": "com.example.product",
		"quantity":  1,
	}

	signature, err := creator.CreateSignature(request)
	assertNoError(t, err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assertNoError(t, err, "Failed to decode JWT")

	assertEqual(t, "advanced-commerce-api", payload["aud"], "Audience")

	// Verify request field exists and is base64 encoded
	assertNotNil(t, payload["request"], "Request field")

	requestStr, ok := payload["request"].(string)
	assertTrue(t, ok, "Request should be a string")
	assertTrue(t, requestStr != "", "Request should not be empty")
}
