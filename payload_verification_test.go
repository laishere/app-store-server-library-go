package appstore

import (
	"testing"
)

// Test app store server notification decoding (environment check)
func TestAppStoreServerNotificationDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assertNoError(t, err, "Failed to read test data")

	notification, err := verifier.VerifyAndDecodeNotification(testNotification)
	assertNoError(t, err, "Failed to verify and decode notification")
	assertEqual(t, NOTIFICATION_TYPE_TEST, notification.NotificationType, "NotificationType")
}

func TestAppStoreServerNotificationDecodingProduction(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", int64Ptr(1234))
	assertNoError(t, err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assertNoError(t, err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	assertError(t, err, "Expected verification to fail with wrong environment")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, INVALID_ENVIRONMENT, verifyErr.Status, "Verification error status")
}

// Test missing x5c header
func TestMissingX5CHeader(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	missingX5CHeaderClaim, err := readTestDataString("mock_signed_data/missingX5CHeaderClaim")
	assertNoError(t, err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(missingX5CHeaderClaim)
	assertError(t, err, "Expected verification to fail with missing x5c")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}

// Test payload verification with wrong bundle ID
func TestWrongBundleIdForServerNotification(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.examplex", nil)
	assertNoError(t, err, "Failed to create verifier")

	wrongBundleData, err := readTestDataString("mock_signed_data/wrongBundleId")
	assertNoError(t, err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(wrongBundleData)
	assertError(t, err, "Expected verification to fail with wrong bundle ID")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, INVALID_APP_IDENTIFIER, verifyErr.Status, "Verification error status")
}

// Test wrong app apple id
func TestWrongAppAppleIdForServerNotification(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", int64Ptr(1235))
	assertNoError(t, err, "Failed to create verifier")

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	assertNoError(t, err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	assertError(t, err, "Expected verification to fail with wrong app apple id")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, INVALID_APP_IDENTIFIER, verifyErr.Status, "Verification error status")
}

func TestRenewalInfoDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	renewalInfoData, err := readTestDataString("mock_signed_data/renewalInfo")
	assertNoError(t, err, "Failed to read test data")

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(renewalInfoData)
	assertNoError(t, err, "Failed to verify and decode renewal info")
	assertEqual(t, ENVIRONMENT_SANDBOX, renewalInfo.Environment, "Environment")
}

func TestTransactionInfoDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	transactionInfoData, err := readTestDataString("mock_signed_data/transactionInfo")
	assertNoError(t, err, "Failed to read test data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(transactionInfoData)
	assertNoError(t, err, "Failed to verify and decode transaction")
	assertEqual(t, ENVIRONMENT_SANDBOX, transaction.Environment, "Environment")
}

// Test malformed JWT
func TestMalformedJwtWithTooManyParts(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	_, err = verifier.VerifyAndDecodeNotification("a.b.c.d")
	assertError(t, err, "Expected verification to fail with malformed JWT")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}

func TestMalformedJwtWithMalformedData(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	assertNoError(t, err, "Failed to create verifier")

	_, err = verifier.VerifyAndDecodeNotification("a.b.c")
	assertError(t, err, "Expected verification to fail with malformed JWT")

	verifyErr, ok := err.(*VerificationException)
	assertTrue(t, ok, "Expected VerificationException")
	assertEqual(t, VERIFICATION_FAILURE, verifyErr.Status, "Verification error status")
}
