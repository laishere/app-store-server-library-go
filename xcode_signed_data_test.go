package appstore

import (
	"testing"
)

const XCODE_BUNDLE_ID = "com.example.naturelab.backyardbirds.example"

// Test Xcode signed app transaction
func TestXcodeSignedAppTransaction(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	encodedAppTransaction, err := readTestDataString("xcode/xcode-signed-app-transaction")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	appTransaction, err := verifier.VerifyAndDecodeAppTransaction(encodedAppTransaction)
	if err != nil {
		t.Fatalf("Failed to verify and decode app transaction: %v", err)
	}

	assertNotNil(t, appTransaction, "AppTransaction")
	assertNil(t, appTransaction.AppAppleId, "AppAppleId")
	assertEqual(t, XCODE_BUNDLE_ID, appTransaction.BundleId, "BundleId")
	assertEqual(t, "1", appTransaction.ApplicationVersion, "ApplicationVersion")
	assertNil(t, appTransaction.VersionExternalIdentifier, "VersionExternalIdentifier")
	assertEqual(t, Timestamp(-62135769600000), appTransaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assertEqual(t, "1", appTransaction.OriginalApplicationVersion, "OriginalApplicationVersion")
	assertEqual(t, "cYUsXc53EbYc0pOeXG5d6/31LGHeVGf84sqSN0OrJi5u/j2H89WWKgS8N0hMsMlf", appTransaction.DeviceVerification, "DeviceVerification")
	assertEqual(t, "48c8b92d-ce0d-4229-bedf-e61b4f9cfc92", appTransaction.DeviceVerificationNonce, "DeviceVerificationNonce")
	assertNil(t, appTransaction.PreorderDate, "PreorderDate")
	assertEqual(t, ENVIRONMENT_XCODE, appTransaction.ReceiptType, "ReceiptType")
}

// Test Xcode signed transaction
func TestXcodeSignedTransaction(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	encodedTransaction, err := readTestDataString("xcode/xcode-signed-transaction")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(encodedTransaction)
	if err != nil {
		t.Fatalf("Failed to verify and decode transaction: %v", err)
	}

	assertEqual(t, "0", transaction.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, "0", transaction.TransactionId, "TransactionId")
	assertEqual(t, "0", transaction.WebOrderLineItemId, "WebOrderLineItemId")
	assertEqual(t, XCODE_BUNDLE_ID, transaction.BundleId, "BundleId")
	assertEqual(t, "pass.premium", transaction.ProductId, "ProductId")
	assertEqual(t, "6F3A93AB", transaction.SubscriptionGroupIdentifier, "SubscriptionGroupIdentifier")
	assertEqual(t, Timestamp(1697679936049), transaction.PurchaseDate, "PurchaseDate")
	assertEqual(t, Timestamp(1697679936049), transaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assertEqual(t, Timestamp(1700358336049), transaction.ExpiresDate, "ExpiresDate")
	assertEqual(t, int32(1), transaction.Quantity, "Quantity")
	assertEqual(t, TYPE_AUTO_RENEWABLE_SUBSCRIPTION, transaction.Type, "Type")
	assertNil(t, transaction.AppAccountToken, "AppAccountToken")
	assertEqual(t, IN_APP_OWNERSHIP_TYPE_PURCHASED, transaction.InAppOwnershipType, "InAppOwnershipType")
	assertEqual(t, Timestamp(1697679936056), transaction.SignedDate, "SignedDate")
	assertNil(t, transaction.RevocationReason, "RevocationReason")
	assertNil(t, transaction.RevocationDate, "RevocationDate")
	assertEqual(t, false, transaction.IsUpgraded, "IsUpgraded")
	assertEqual(t, OFFER_TYPE_INTRODUCTORY, transaction.OfferType, "OfferType")
	assertNil(t, transaction.OfferIdentifier, "OfferIdentifier")
	assertEqual(t, ENVIRONMENT_XCODE, transaction.Environment, "Environment")
	assertEqual(t, "USA", transaction.Storefront, "Storefront")
	assertEqual(t, "143441", transaction.StorefrontId, "StorefrontId")
	assertEqual(t, TRANSACTION_REASON_PURCHASE, transaction.TransactionReason, "TransactionReason")
}

// Test Xcode signed renewal info
func TestXcodeSignedRenewalInfo(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	encodedRenewalInfo, err := readTestDataString("xcode/xcode-signed-renewal-info")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(encodedRenewalInfo)
	if err != nil {
		t.Fatalf("Failed to verify and decode renewal info: %v", err)
	}

	assertNil(t, renewalInfo.ExpirationIntent, "ExpirationIntent")
	assertEqual(t, "0", renewalInfo.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, "pass.premium", renewalInfo.AutoRenewProductId, "AutoRenewProductId")
	assertEqual(t, "pass.premium", renewalInfo.ProductId, "ProductId")
	assertEqual(t, AUTO_RENEW_STATUS_ON, renewalInfo.AutoRenewStatus, "AutoRenewStatus")
	assertNil(t, renewalInfo.IsInBillingRetryPeriod, "IsInBillingRetryPeriod")
	assertNil(t, renewalInfo.PriceIncreaseStatus, "PriceIncreaseStatus")
	assertNil(t, renewalInfo.GracePeriodExpiresDate, "GracePeriodExpiresDate")
	assertNil(t, renewalInfo.OfferType, "OfferType")
	assertNil(t, renewalInfo.OfferIdentifier, "OfferIdentifier")
	assertEqual(t, Timestamp(1697679936711), renewalInfo.SignedDate, "SignedDate")
	assertEqual(t, ENVIRONMENT_XCODE, renewalInfo.Environment, "Environment")
	assertEqual(t, Timestamp(1697679936049), renewalInfo.RecentSubscriptionStartDate, "RecentSubscriptionStartDate")
	assertEqual(t, Timestamp(1700358336049), renewalInfo.RenewalDate, "RenewalDate")
}

// Test Xcode signed app transaction with production environment should fail
func TestXcodeSignedAppTransactionWithProductionEnvironment(t *testing.T) {
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, XCODE_BUNDLE_ID, int64Ptr(1234))
	if err != nil {
		t.Fatalf("Failed to create verifier: %v", err)
	}

	encodedAppTransaction, err := readTestDataString("xcode/xcode-signed-app-transaction")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	_, err = verifier.VerifyAndDecodeAppTransaction(encodedAppTransaction)
	if err == nil {
		t.Fatal("Expected verification to fail with wrong environment")
	}

	// Just verify it's a VerificationException - the specific status may vary
	// depending on whether chain verification fails first or environment check fails first
	_, ok := err.(*VerificationException)
	if !ok {
		t.Fatalf("Expected VerificationException, got %T", err)
	}
}
