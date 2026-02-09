package appstore

import (
	"testing"
)

// Test app store server notification decoding (environment check)
func TestAppStoreServerNotificationDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	notification, err := verifier.VerifyAndDecodeNotification(testNotification)
	if err != nil {
		t.Fatalf("Failed to verify and decode notification: %v", err)
	}
	assertEqual(t, NOTIFICATION_TYPE_TEST, notification.NotificationType, "NotificationType")
}

func TestAppStoreServerNotificationDecodingProduction(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", int64Ptr(1234))
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	if err == nil {
		t.Fatal("Expected verification to fail with wrong environment")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != INVALID_ENVIRONMENT {
		t.Fatalf("Expected INVALID_ENVIRONMENT status, got %v", err)
	}
}

// Test missing x5c header
func TestMissingX5CHeader(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	missingX5CHeaderClaim, err := readTestDataString("mock_signed_data/missingX5CHeaderClaim")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification(missingX5CHeaderClaim)
	if err == nil {
		t.Fatal("Expected verification to fail with missing x5c")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != VERIFICATION_FAILURE {
		t.Fatalf("Expected VERIFICATION_FAILURE status, got %v", err)
	}
}

// Test payload verification with wrong bundle ID
func TestWrongBundleIdForServerNotification(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.examplex", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	wrongBundleData, err := readTestDataString("mock_signed_data/wrongBundleId")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification(wrongBundleData)
	if err == nil {
		t.Fatal("Expected verification to fail with wrong bundle ID")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != INVALID_APP_IDENTIFIER {
		t.Fatalf("Expected INVALID_APP_IDENTIFIER status, got %v", err)
	}
}

// Test wrong app apple id
func TestWrongAppAppleIdForServerNotification(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, "com.example", int64Ptr(1235))
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	testNotification, err := readTestDataString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	if err == nil {
		t.Fatal("Expected verification to fail with wrong app apple id")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != INVALID_APP_IDENTIFIER {
		t.Fatalf("Expected INVALID_APP_IDENTIFIER status, got %v", err)
	}
}

func TestRenewalInfoDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	renewalInfoData, err := readTestDataString("mock_signed_data/renewalInfo")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(renewalInfoData)
	if err != nil {
		t.Fatalf("Failed to verify and decode renewal info: %v", err)
	}
	assertEqual(t, ENVIRONMENT_SANDBOX, renewalInfo.Environment, "Environment")
}

func TestTransactionInfoDecoding(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	transactionInfoData, err := readTestDataString("mock_signed_data/transactionInfo")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(transactionInfoData)
	if err != nil {
		t.Fatalf("Failed to verify and decode transaction: %v", err)
	}
	assertEqual(t, ENVIRONMENT_SANDBOX, transaction.Environment, "Environment")
}

// Test malformed JWT
func TestMalformedJwtWithTooManyParts(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification("a.b.c.d")
	if err == nil {
		t.Fatal("Expected verification to fail with malformed JWT")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != VERIFICATION_FAILURE {
		t.Fatalf("Expected VERIFICATION_FAILURE status, got %v", err)
	}
}

func TestMalformedJwtWithMalformedData(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_SANDBOX, "com.example", nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	_, err = verifier.VerifyAndDecodeNotification("a.b.c")
	if err == nil {
		t.Fatal("Expected verification to fail with malformed JWT")
	}

	verifyErr, ok := err.(*VerificationException)
	if !ok || verifyErr.Status != VERIFICATION_FAILURE {
		t.Fatalf("Expected VERIFICATION_FAILURE status, got %v", err)
	}
}
