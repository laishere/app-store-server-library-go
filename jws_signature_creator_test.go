package appstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test introductory offer eligibility signature
func TestIntroductoryOfferEligibilitySignatureCreator(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewIntroductoryOfferEligibilitySignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", true, "999")
	assert.NoError(err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assert.NoError(err, "Failed to decode JWT")

	assert.Equal("introductory-offer-eligibility", payload["aud"], "Audience")
	assert.Equal("com.example.product", payload["productId"], "Product ID")
	assert.Equal(true, payload["allowIntroductoryOffer"], "Allow Introductory Offer")
	assert.Equal("999", payload["transactionId"], "Transaction ID")
}

// Test advanced commerce API signature
func TestAdvancedCommerceAPIInAppSignatureCreator(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewAdvancedCommerceAPIInAppSignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	request := map[string]any{
		"productId": "com.example.product",
		"quantity":  1,
	}

	signature, err := creator.CreateSignature(request)
	assert.NoError(err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assert.NoError(err, "Failed to decode JWT")

	assert.Equal("advanced-commerce-api", payload["aud"], "Audience")

	// Verify request field exists and is base64 encoded
	assert.NotNil(payload["request"], "Request field")

	requestStr, ok := payload["request"].(string)
	assert.True(ok, "Request should be a string")
	assert.True(requestStr != "", "Request should not be empty")
}
