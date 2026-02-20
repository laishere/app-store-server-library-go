package appstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test promotional offer V2 signature creation
func TestPromotionalOfferV2SignatureCreator(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", ptr("999"))
	assert.NoError(err, "Failed to create signature")

	// Verify the JWT was created
	assert.True(signature != "", "Signature is not empty")

	// Parse and verify the JWT structure
	header, payload, err := decodeJWTWithoutVerification(signature)
	assert.NoError(err, "Failed to decode JWT")

	// Verify header
	assert.Equal("ES256", header["alg"], "Algorithm")
	assert.Equal(TEST_KEY_ID, header["kid"], "Key ID")

	// Verify payload
	assert.Equal(TEST_BUNDLE_ID, payload["bid"], "Bundle ID")
	assert.Equal(TEST_ISSUER_ID, payload["iss"], "Issuer ID")
	assert.Equal("promotional-offer", payload["aud"], "Audience")
	assert.Equal("com.example.product", payload["productId"], "Product ID")
	assert.Equal("OFFER123", payload["offerIdentifier"], "Offer Identifier")
	assert.Equal("999", payload["transactionId"], "Transaction ID")

	// Verify nonce exists
	assert.True(payload["nonce"] != nil && payload["nonce"] != "", "Nonce is present")

	// Verify iat exists
	assert.NotNil(payload["iat"], "IAT")
}

// Test promotional offer V2 without transaction ID
func TestPromotionalOfferV2SignatureCreator_NoTransactionId(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	signature, err := creator.CreateSignature("com.example.product", "OFFER123", nil)
	assert.NoError(err, "Failed to create signature")

	_, payload, err := decodeJWTWithoutVerification(signature)
	assert.NoError(err, "Failed to decode JWT")

	// Transaction ID should not be in payload
	assert.Nil(payload["transactionId"], "Transaction ID")
}

// Test promotional offer V2 with missing required fields
func TestPromotionalOfferV2SignatureCreator_MissingProductId(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	_, err = creator.CreateSignature("", "OFFER123", nil)
	assert.Error(err, "Expected error for empty product ID")
}

func TestPromotionalOfferV2SignatureCreator_MissingOfferIdentifier(t *testing.T) {
	assert := assert.New(t)
	keyBytes, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	creator, err := NewPromotionalOfferV2SignatureCreator(keyBytes, TEST_KEY_ID, TEST_ISSUER_ID, TEST_BUNDLE_ID)
	assert.NoError(err, "Failed to create signature creator")

	_, err = creator.CreateSignature("com.example.product", "", nil)
	assert.Error(err, "Expected error for empty offer identifier")
}
