package appstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test app store server notification decoding (environment check)
func TestAppStoreServerNotificationDecoding(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assert.NoError(err, "Failed to read test data")

	notification, err := verifier.VerifyAndDecodeNotification(testNotification)
	assert.NoError(err, "Failed to verify and decode notification")
	assert.Equal(NOTIFICATION_TYPE_TEST, notification.NotificationType, "NotificationType")
}

func TestAppStoreServerNotificationDecodingProduction(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", ptr(int64(1234)))
	assert.NoError(err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assert.NoError(err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	assert.Error(err, "Expected verification to fail with wrong environment")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(INVALID_ENVIRONMENT, verifyErr.Status, "Verification error status")
}

// Test missing x5c header
func TestMissingX5CHeader(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	missingX5CHeaderClaim, err := readTestDataString("mock_signed_data/missingX5CHeaderClaim")
	assert.NoError(err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(missingX5CHeaderClaim)
	assert.Error(err, "Expected verification to fail with missing x5c")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}

// Test payload verification with wrong bundle ID
func TestWrongBundleIdForServerNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.examplex", nil)
	assert.NoError(err, "Failed to create verifier")

	wrongBundleData, err := readTestDataString("mock_signed_data/wrongBundleId")
	assert.NoError(err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(wrongBundleData)
	assert.Error(err, "Expected verification to fail with wrong bundle ID")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(INVALID_APP_IDENTIFIER, verifyErr.Status, "Verification error status")
}

// Test wrong app apple id
func TestWrongAppAppleIdForServerNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", ptr(int64(1235)))
	assert.NoError(err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assert.NoError(err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	assert.Error(err, "Expected verification to fail with wrong app apple id")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(INVALID_APP_IDENTIFIER, verifyErr.Status, "Verification error status")
}

func TestRenewalInfoDecoding(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	renewalInfoData, err := readTestDataString("mock_signed_data/renewalInfo")
	assert.NoError(err, "Failed to read test data")

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(renewalInfoData)
	assert.NoError(err, "Failed to verify and decode renewal info")
	assert.Equal(ENVIRONMENT_SANDBOX, renewalInfo.Environment, "Environment")
}

func TestTransactionInfoDecoding(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	transactionInfoData, err := readTestDataString("mock_signed_data/transactionInfo")
	assert.NoError(err, "Failed to read test data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(transactionInfoData)
	assert.NoError(err, "Failed to verify and decode transaction")
	assert.Equal(ENVIRONMENT_SANDBOX, transaction.Environment, "Environment")
}

// Test malformed JWT
func TestMalformedJwtWithTooManyParts(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	_, err = verifier.VerifyAndDecodeNotification("a.b.c.d")
	assert.Error(err, "Expected verification to fail with malformed JWT")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}

func TestMalformedJwtWithMalformedData(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assert.NoError(err, "Failed to create verifier")

	_, err = verifier.VerifyAndDecodeNotification("a.b.c")
	assert.Error(err, "Expected verification to fail with malformed JWT")

	verifyErr, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
	assert.Equal(VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}
