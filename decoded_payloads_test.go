package appstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// createDefaultTestSignedDataVerifier creates a verifier with default test settings
func createDefaultTestSignedDataVerifier() (*SignedDataVerifier, error) {
	return createTestSignedDataVerifier(ENVIRONMENT_LOCAL_TESTING, "com.example", nil)
}

// Test app transaction decoding
func TestVerifyAndDecodeAppTransaction(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedAppTransaction, err := createSignedDataFromJSON("models/appTransaction.json")
	assert.NoError(err, "Failed to create signed data")

	appTransaction, err := verifier.VerifyAndDecodeAppTransaction(signedAppTransaction)
	assert.NoError(err, "Failed to verify and decode app transaction")

	// Verify fields
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, appTransaction.ReceiptType, "ReceiptType")
	assert.Equal(int64(531412), *appTransaction.AppAppleId, "AppAppleId")
	assert.Equal("com.example", appTransaction.BundleId, "BundleId")
	assert.Equal("1.2.3", appTransaction.ApplicationVersion, "ApplicationVersion")
	assert.Equal(int64(512), *appTransaction.VersionExternalIdentifier, "VersionExternalIdentifier")
	assert.Equal(Timestamp(1698148900000), appTransaction.ReceiptCreationDate, "ReceiptCreationDate")
	assert.Equal(Timestamp(1698148800000), appTransaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assert.Equal("1.1.2", appTransaction.OriginalApplicationVersion, "OriginalApplicationVersion")
	assert.Equal("device_verification_value", appTransaction.DeviceVerification, "DeviceVerification")
	assert.Equal("48ccfa42-7431-4f22-9908-7e88983e105a", appTransaction.DeviceVerificationNonce, "DeviceVerificationNonce")
	assert.Equal(Timestamp(1698148700000), *appTransaction.PreorderDate, "PreorderDate")
}

// Test transaction decoding
func TestVerifyAndDecodeSignedTransaction(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedTransaction, err := createSignedDataFromJSON("models/signedTransaction.json")
	assert.NoError(err, "Failed to create signed data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	assert.NoError(err, "Failed to verify and decode transaction")

	assert.Equal("12345", transaction.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal("23456", transaction.TransactionId, "TransactionId")
	assert.Equal("34343", transaction.WebOrderLineItemId, "WebOrderLineItemId")
	assert.Equal("com.example", transaction.BundleId, "BundleId")
	assert.Equal("com.example.product", transaction.ProductId, "ProductId")
	assert.Equal("55555", transaction.SubscriptionGroupIdentifier, "SubscriptionGroupIdentifier")
	assert.Equal(Timestamp(1698148800000), transaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assert.Equal(Timestamp(1698148900000), transaction.PurchaseDate, "PurchaseDate")
	assert.Equal(Timestamp(1698148950000), *transaction.RevocationDate, "RevocationDate")
	assert.Equal(Timestamp(1698149000000), transaction.ExpiresDate, "ExpiresDate")
	assert.Equal(int32(1), transaction.Quantity, "Quantity")
	assert.Equal(TYPE_AUTO_RENEWABLE_SUBSCRIPTION, transaction.Type, "Type")
	assert.Equal("7e3fb20b-4cdb-47cc-936d-99d65f608138", *transaction.AppAccountToken, "AppAccountToken")
	assert.Equal(IN_APP_OWNERSHIP_TYPE_PURCHASED, transaction.InAppOwnershipType, "InAppOwnershipType")
	assert.Equal(Timestamp(1698148900000), transaction.SignedDate, "SignedDate")
	assert.Equal(REVOCATION_REASON_REFUNDED_DUE_TO_ISSUE, *transaction.RevocationReason, "RevocationReason")
	assert.Equal("abc.123", *transaction.OfferIdentifier, "OfferIdentifier")
	assert.Equal(true, transaction.IsUpgraded, "IsUpgraded")
	assert.Equal(OFFER_TYPE_INTRODUCTORY, transaction.OfferType, "OfferType")
	assert.Equal("USA", transaction.Storefront, "Storefront")
	assert.Equal("143441", transaction.StorefrontId, "StorefrontId")
	assert.Equal(TRANSACTION_REASON_PURCHASE, transaction.TransactionReason, "TransactionReason")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, transaction.Environment, "Environment")
	assert.Equal(int64(10990), transaction.Price, "Price")
	assert.Equal("USD", transaction.Currency, "Currency")
	assert.Equal(OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO, transaction.OfferDiscountType, "OfferDiscountType")
}

// Test transaction with revocation
func TestVerifyAndDecodeSignedTransactionWithRevocation(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedTransaction, err := createSignedDataFromJSON("models/signedTransactionWithRevocation.json")
	assert.NoError(err, "Failed to create signed data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	assert.NoError(err, "Failed to verify and decode transaction")

	// Verify revocation fields
	assert.Equal(REVOCATION_TYPE_REFUND_PRORATED, transaction.RevocationType, "RevocationType")
	assert.Equal(int32(50000), transaction.RevocationPercentage, "RevocationPercentage")
}

// Test renewal info decoding
func TestVerifyAndDecodeRenewalInfo(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedRenewalInfo, err := createSignedDataFromJSON("models/signedRenewalInfo.json")
	assert.NoError(err, "Failed to create signed data")

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(signedRenewalInfo)
	assert.NoError(err, "Failed to verify and decode renewal info")

	// Verify all fields
	assert.Equal(EXPIRATION_INTENT_CUSTOMER_CANCELLED, *renewalInfo.ExpirationIntent, "ExpirationIntent")
	assert.Equal("12345", renewalInfo.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal("com.example.product.2", renewalInfo.AutoRenewProductId, "AutoRenewProductId")
	assert.Equal("com.example.product", renewalInfo.ProductId, "ProductId")
	assert.Equal(AUTO_RENEW_STATUS_ON, renewalInfo.AutoRenewStatus, "AutoRenewStatus")
	assert.Equal(true, *renewalInfo.IsInBillingRetryPeriod, "IsInBillingRetryPeriod")
	assert.Equal(PRICE_INCREASE_STATUS_CUSTOMER_HAS_NOT_RESPONDED, *renewalInfo.PriceIncreaseStatus, "PriceIncreaseStatus")
	assert.Equal(Timestamp(1698148900000), *renewalInfo.GracePeriodExpiresDate, "GracePeriodExpiresDate")
	assert.Equal(OFFER_TYPE_PROMOTIONAL, *renewalInfo.OfferType, "OfferType")
	assert.Equal("abc.123", *renewalInfo.OfferIdentifier, "OfferIdentifier")
	assert.Equal(Timestamp(1698148800000), renewalInfo.SignedDate, "SignedDate")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, renewalInfo.Environment, "Environment")
	assert.Equal(Timestamp(1698148800000), renewalInfo.RecentSubscriptionStartDate, "RecentSubscriptionStartDate")
	assert.Equal(Timestamp(1698148850000), renewalInfo.RenewalDate, "RenewalDate")
	assert.Equal(int64(9990), renewalInfo.RenewalPrice, "RenewalPrice")
	assert.Equal("USD", renewalInfo.Currency, "Currency")
	assert.Equal(OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO, renewalInfo.OfferDiscountType, "OfferDiscountType")
}

// Test notification decoding (SUBSCRIBED type)
func TestVerifyAndDecodeNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	// Verify notification fields
	assert.Equal(NOTIFICATION_TYPE_SUBSCRIBED, notification.NotificationType, "NotificationType")
	assert.NotNil(notification.Subtype, "Subtype")
	assert.Equal(SUBTYPE_INITIAL_BUY, *notification.Subtype, "Subtype")
	assert.Equal("002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assert.Equal("2.0", notification.Version, "Version")
	assert.Equal(Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assert.NotNil(notification.Data, "Data")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, notification.Data.Environment, "Data.Environment")
	assert.Equal(int64(41234), notification.Data.AppAppleId, "Data.AppAppleId")
	assert.Equal("com.example", notification.Data.BundleId, "Data.BundleId")
	assert.Equal("1.2.3", notification.Data.BundleVersion, "Data.BundleVersion")
	assert.Equal("signed_transaction_info_value", notification.Data.SignedTransactionInfo, "Data.SignedTransactionInfo")
	assert.Equal("signed_renewal_info_value", notification.Data.SignedRenewalInfo, "Data.SignedRenewalInfo")
	assert.Equal(STATUS_ACTIVE, notification.Data.Status, "Data.Status")

	// Verify nil fields
	assert.Nil(notification.Summary, "Summary")
	assert.Nil(notification.ExternalPurchaseToken, "ExternalPurchaseToken")
}

// Test consumption request notification
func TestVerifyAndDecodeConsumptionRequestNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedConsumptionRequestNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	assert.Equal(NOTIFICATION_TYPE_CONSUMPTION_REQUEST, notification.NotificationType, "NotificationType")

	assert.NotNil(notification.Data, "Data")
	assert.Equal(CONSUMPTION_REQUEST_REASON_UNINTENDED_PURCHASE, *notification.Data.ConsumptionRequestReason, "ConsumptionRequestReason")
}

// Test summary notification
func TestVerifyAndDecodeSummaryNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedSummaryNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	assert.Equal(NOTIFICATION_TYPE_RENEWAL_EXTENSION, notification.NotificationType, "NotificationType")
	assert.NotNil(notification.Subtype, "Subtype")
	assert.Equal(SUBTYPE_SUMMARY, *notification.Subtype, "Subtype")

	assert.Nil(notification.Data, "Data")
	assert.NotNil(notification.Summary, "Summary")

	assert.Equal(ENVIRONMENT_LOCAL_TESTING, notification.Summary.Environment, "Summary.Environment")
	assert.Equal(int64(41234), notification.Summary.AppAppleId, "Summary.AppAppleId")
	assert.Equal("com.example", notification.Summary.BundleId, "Summary.BundleId")
	assert.Equal("com.example.product", notification.Summary.ProductId, "Summary.ProductId")
	assert.Equal(int64(5), notification.Summary.SucceededCount, "Summary.SucceededCount")
	assert.Equal(int64(2), notification.Summary.FailedCount, "Summary.FailedCount")
}

func TestVerifyAndDecodeExternalPurchaseTokenNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")
	verifier.allowAnyEnvironment = true

	signedExternalPurchaseTokenNotification, err := createSignedDataFromJSON("models/signedExternalPurchaseTokenNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedExternalPurchaseTokenNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	assert.Equal(NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION, notification.NotificationType, "NotificationType")
	assert.NotNil(notification.Subtype, "Subtype")
	assert.Equal(SUBTYPE_UNREPORTED, *notification.Subtype, "Subtype")
	assert.Equal("002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assert.Equal("2.0", notification.Version, "Version")
	assert.Equal(Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assert.Nil(notification.Data, "Data")
	assert.Nil(notification.Summary, "Summary")
	assert.NotNil(notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	assert.Equal("b2158121-7af9-49d4-9561-1f588205523e", notification.ExternalPurchaseToken.ExternalPurchaseId, "ExternalPurchaseToken.ExternalPurchaseId")
	assert.Equal(Timestamp(1698148950000), notification.ExternalPurchaseToken.TokenCreationDate, "ExternalPurchaseToken.TokenCreationDate")
	assert.Equal(int64(55555), notification.ExternalPurchaseToken.AppAppleId, "ExternalPurchaseToken.AppAppleId")
	assert.Equal("com.example", notification.ExternalPurchaseToken.BundleId, "ExternalPurchaseToken.BundleId")
}

func TestVerifyAndDecodeExternalPurchaseTokenSandboxNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")
	verifier.allowAnyEnvironment = true

	signedExternalPurchaseTokenNotification, err := createSignedDataFromJSON("models/signedExternalPurchaseTokenSandboxNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedExternalPurchaseTokenNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	assert.Equal(NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION, notification.NotificationType, "NotificationType")
	assert.NotNil(notification.Subtype, "Subtype")
	assert.Equal(SUBTYPE_UNREPORTED, *notification.Subtype, "Subtype")
	assert.Equal("002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assert.Equal("2.0", notification.Version, "Version")
	assert.Equal(Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assert.Nil(notification.Data, "Data")
	assert.Nil(notification.Summary, "Summary")
	assert.NotNil(notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	assert.Equal("SANDBOX_b2158121-7af9-49d4-9561-1f588205523e", notification.ExternalPurchaseToken.ExternalPurchaseId, "ExternalPurchaseToken.ExternalPurchaseId")
	assert.Equal(Timestamp(1698148950000), notification.ExternalPurchaseToken.TokenCreationDate, "ExternalPurchaseToken.TokenCreationDate")
	assert.Equal(int64(55555), notification.ExternalPurchaseToken.AppAppleId, "ExternalPurchaseToken.AppAppleId")
	assert.Equal("com.example", notification.ExternalPurchaseToken.BundleId, "ExternalPurchaseToken.BundleId")
}

// Test rescind consent notification
func TestVerifyAndDecodeRescindConsentNotification(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedRescindConsentNotification.json")
	assert.NoError(err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assert.NoError(err, "Failed to verify and decode notification")

	// Verify notification type
	assert.Equal(NOTIFICATION_TYPE_RESCIND_CONSENT, notification.NotificationType, "NotificationType")
	assert.Equal("002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assert.Equal("2.0", notification.Version, "Version")
	assert.Equal(Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	// Verify nil fields
	assert.Nil(notification.Subtype, "Subtype")
	assert.Nil(notification.Data, "Data")
	assert.Nil(notification.Summary, "Summary")
	assert.Nil(notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	// Verify AppData - THIS TESTS AppData unmarshal!
	assert.NotNil(notification.AppData, "AppData")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, notification.AppData.Environment, "AppData.Environment")
	assert.Equal(int64(41234), notification.AppData.AppAppleId, "AppData.AppAppleId")
	assert.Equal("com.example", notification.AppData.BundleId, "AppData.BundleId")
	assert.Equal("signed_app_transaction_info_value", notification.AppData.SignedAppTransactionInfo, "AppData.SignedAppTransactionInfo")
}

// Test realtime request decoding
func TestVerifyAndDecodeRealtimeRequest(t *testing.T) {
	assert := assert.New(t)
	verifier, err := createDefaultTestSignedDataVerifier()
	assert.NoError(err, "Failed to create verifier")

	signedRealtimeRequest, err := createSignedDataFromJSON("models/decodedRealtimeRequest.json")
	assert.NoError(err, "Failed to create signed data")

	request, err := verifier.VerifyAndDecodeRealtimeRequest(signedRealtimeRequest)
	assert.NoError(err, "Failed to verify and decode realtime request")

	// Verify fields - THIS TESTS DecodedRealtimeRequestBody unmarshal!
	assert.Equal("99371282", request.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal(int64(531412), request.AppAppleId, "AppAppleId")
	assert.Equal("com.example.product", request.ProductId, "ProductId")
	assert.Equal("en-US", request.UserLocale, "UserLocale")
	assert.Equal("3db5c98d-8acf-4e29-831e-8e1f82f9f6e9", request.RequestIdentifier, "RequestIdentifier")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, request.Environment, "Environment")
	assert.Equal(Timestamp(1698148900000), request.SignedDate, "SignedDate")
}
