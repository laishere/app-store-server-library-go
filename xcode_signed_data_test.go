package appstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const XCODE_BUNDLE_ID = "com.example.naturelab.backyardbirds.example"

// Test Xcode signed app transaction
func TestXcodeSignedAppTransaction(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	assert.NoError(err, "Failed to create verifier")

	encodedAppTransaction, err := readTestDataString("xcode/xcode-signed-app-transaction")
	assert.NoError(err, "Failed to read test data")

	appTransaction, err := verifier.VerifyAndDecodeAppTransaction(encodedAppTransaction)
	assert.NoError(err, "Failed to verify and decode app transaction")

	assert.NotNil(appTransaction, "AppTransaction")
	assert.Nil(appTransaction.AppAppleId, "AppAppleId")
	assert.Equal(XCODE_BUNDLE_ID, appTransaction.BundleId, "BundleId")
	assert.Equal("1", appTransaction.ApplicationVersion, "ApplicationVersion")
	assert.Nil(appTransaction.VersionExternalIdentifier, "VersionExternalIdentifier")
	assert.Equal(Timestamp(-62135769600000), appTransaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assert.Equal("1", appTransaction.OriginalApplicationVersion, "OriginalApplicationVersion")
	assert.Equal("cYUsXc53EbYc0pOeXG5d6/31LGHeVGf84sqSN0OrJi5u/j2H89WWKgS8N0hMsMlf", appTransaction.DeviceVerification, "DeviceVerification")
	assert.Equal("48c8b92d-ce0d-4229-bedf-e61b4f9cfc92", appTransaction.DeviceVerificationNonce, "DeviceVerificationNonce")
	assert.Nil(appTransaction.PreorderDate, "PreorderDate")
	assert.Equal(ENVIRONMENT_XCODE, appTransaction.ReceiptType, "ReceiptType")
}

// Test Xcode signed transaction
func TestXcodeSignedTransaction(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	assert.NoError(err, "Failed to create verifier")

	encodedTransaction, err := readTestDataString("xcode/xcode-signed-transaction")
	assert.NoError(err, "Failed to read test data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(encodedTransaction)
	assert.NoError(err, "Failed to verify and decode transaction")

	assert.Equal("0", transaction.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal("0", transaction.TransactionId, "TransactionId")
	assert.Equal("0", transaction.WebOrderLineItemId, "WebOrderLineItemId")
	assert.Equal(XCODE_BUNDLE_ID, transaction.BundleId, "BundleId")
	assert.Equal("pass.premium", transaction.ProductId, "ProductId")
	assert.Equal("6F3A93AB", transaction.SubscriptionGroupIdentifier, "SubscriptionGroupIdentifier")
	assert.Equal(Timestamp(1697679936049), transaction.PurchaseDate, "PurchaseDate")
	assert.Equal(Timestamp(1697679936049), transaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assert.Equal(Timestamp(1700358336049), transaction.ExpiresDate, "ExpiresDate")
	assert.Equal(int32(1), transaction.Quantity, "Quantity")
	assert.Equal(TYPE_AUTO_RENEWABLE_SUBSCRIPTION, transaction.Type, "Type")
	assert.Nil(transaction.AppAccountToken, "AppAccountToken")
	assert.Equal(IN_APP_OWNERSHIP_TYPE_PURCHASED, transaction.InAppOwnershipType, "InAppOwnershipType")
	assert.Equal(Timestamp(1697679936056), transaction.SignedDate, "SignedDate")
	assert.Nil(transaction.RevocationReason, "RevocationReason")
	assert.Nil(transaction.RevocationDate, "RevocationDate")
	assert.Equal(false, transaction.IsUpgraded, "IsUpgraded")
	assert.Equal(OFFER_TYPE_INTRODUCTORY, transaction.OfferType, "OfferType")
	assert.Nil(transaction.OfferIdentifier, "OfferIdentifier")
	assert.Equal(ENVIRONMENT_XCODE, transaction.Environment, "Environment")
	assert.Equal("USA", transaction.Storefront, "Storefront")
	assert.Equal("143441", transaction.StorefrontId, "StorefrontId")
	assert.Equal(TRANSACTION_REASON_PURCHASE, transaction.TransactionReason, "TransactionReason")
}

// Test Xcode signed renewal info
func TestXcodeSignedRenewalInfo(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_XCODE, XCODE_BUNDLE_ID, nil)
	assert.NoError(err, "Failed to create verifier")

	encodedRenewalInfo, err := readTestDataString("xcode/xcode-signed-renewal-info")
	assert.NoError(err, "Failed to read test data")

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(encodedRenewalInfo)
	assert.NoError(err, "Failed to verify and decode renewal info")

	assert.Nil(renewalInfo.ExpirationIntent, "ExpirationIntent")
	assert.Equal("0", renewalInfo.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal("pass.premium", renewalInfo.AutoRenewProductId, "AutoRenewProductId")
	assert.Equal("pass.premium", renewalInfo.ProductId, "ProductId")
	assert.Equal(AUTO_RENEW_STATUS_ON, renewalInfo.AutoRenewStatus, "AutoRenewStatus")
	assert.Nil(renewalInfo.IsInBillingRetryPeriod, "IsInBillingRetryPeriod")
	assert.Nil(renewalInfo.PriceIncreaseStatus, "PriceIncreaseStatus")
	assert.Nil(renewalInfo.GracePeriodExpiresDate, "GracePeriodExpiresDate")
	assert.Nil(renewalInfo.OfferType, "OfferType")
	assert.Nil(renewalInfo.OfferIdentifier, "OfferIdentifier")
	assert.Equal(Timestamp(1697679936711), renewalInfo.SignedDate, "SignedDate")
	assert.Equal(ENVIRONMENT_XCODE, renewalInfo.Environment, "Environment")
	assert.Equal(Timestamp(1697679936049), renewalInfo.RecentSubscriptionStartDate, "RecentSubscriptionStartDate")
	assert.Equal(Timestamp(1700358336049), renewalInfo.RenewalDate, "RenewalDate")
}

// Test Xcode signed app transaction with production environment should fail
func TestXcodeSignedAppTransactionWithProductionEnvironment(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createTestSignedDataVerifier(ENVIRONMENT_PRODUCTION, XCODE_BUNDLE_ID, ptr(int64(1234)))
	assert.NoError(err, "Failed to create verifier")

	encodedAppTransaction, err := readTestDataString("xcode/xcode-signed-app-transaction")
	assert.NoError(err, "Failed to read test data")

	_, err = verifier.VerifyAndDecodeAppTransaction(encodedAppTransaction)
	assert.Error(err, "Expected verification to fail with wrong environment")

	// Just verify it's a VerificationException - the specific status may vary
	// depending on whether chain verification fails first or environment check fails first
	_, ok := err.(*VerificationException)
	assert.True(ok, "Expected VerificationException")
}
