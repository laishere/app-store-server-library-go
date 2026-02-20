package appstore

import (
	"testing"
)

// createDefaultTestSignedDataVerifier creates a verifier with default test settings
func createDefaultTestSignedDataVerifier() (*SignedDataVerifier, error) {
	return createTestSignedDataVerifier(ENVIRONMENT_LOCAL_TESTING, "com.example", nil)
}

// Test app transaction decoding
func TestVerifyAndDecodeAppTransaction(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedAppTransaction, err := createSignedDataFromJSON("models/appTransaction.json")
	assertNoError(t, err, "Failed to create signed data")

	appTransaction, err := verifier.VerifyAndDecodeAppTransaction(signedAppTransaction)
	assertNoError(t, err, "Failed to verify and decode app transaction")

	// Verify fields
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, appTransaction.ReceiptType, "ReceiptType")
	assertEqual(t, int64(531412), *appTransaction.AppAppleId, "AppAppleId")
	assertEqual(t, "com.example", appTransaction.BundleId, "BundleId")
	assertEqual(t, "1.2.3", appTransaction.ApplicationVersion, "ApplicationVersion")
	assertEqual(t, int64(512), *appTransaction.VersionExternalIdentifier, "VersionExternalIdentifier")
	assertEqual(t, Timestamp(1698148900000), appTransaction.ReceiptCreationDate, "ReceiptCreationDate")
	assertEqual(t, Timestamp(1698148800000), appTransaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assertEqual(t, "1.1.2", appTransaction.OriginalApplicationVersion, "OriginalApplicationVersion")
	assertEqual(t, "device_verification_value", appTransaction.DeviceVerification, "DeviceVerification")
	assertEqual(t, "48ccfa42-7431-4f22-9908-7e88983e105a", appTransaction.DeviceVerificationNonce, "DeviceVerificationNonce")
	assertEqual(t, Timestamp(1698148700000), *appTransaction.PreorderDate, "PreorderDate")
}

// Test transaction decoding
func TestVerifyAndDecodeSignedTransaction(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedTransaction, err := createSignedDataFromJSON("models/signedTransaction.json")
	assertNoError(t, err, "Failed to create signed data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	assertNoError(t, err, "Failed to verify and decode transaction")

	assertEqual(t, "12345", transaction.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, "23456", transaction.TransactionId, "TransactionId")
	assertEqual(t, "34343", transaction.WebOrderLineItemId, "WebOrderLineItemId")
	assertEqual(t, "com.example", transaction.BundleId, "BundleId")
	assertEqual(t, "com.example.product", transaction.ProductId, "ProductId")
	assertEqual(t, "55555", transaction.SubscriptionGroupIdentifier, "SubscriptionGroupIdentifier")
	assertEqual(t, Timestamp(1698148800000), transaction.OriginalPurchaseDate, "OriginalPurchaseDate")
	assertEqual(t, Timestamp(1698148900000), transaction.PurchaseDate, "PurchaseDate")
	assertEqual(t, Timestamp(1698148950000), *transaction.RevocationDate, "RevocationDate")
	assertEqual(t, Timestamp(1698149000000), transaction.ExpiresDate, "ExpiresDate")
	assertEqual(t, int32(1), transaction.Quantity, "Quantity")
	assertEqual(t, TYPE_AUTO_RENEWABLE_SUBSCRIPTION, transaction.Type, "Type")
	assertEqual(t, "7e3fb20b-4cdb-47cc-936d-99d65f608138", *transaction.AppAccountToken, "AppAccountToken")
	assertEqual(t, IN_APP_OWNERSHIP_TYPE_PURCHASED, transaction.InAppOwnershipType, "InAppOwnershipType")
	assertEqual(t, Timestamp(1698148900000), transaction.SignedDate, "SignedDate")
	assertEqual(t, REVOCATION_REASON_REFUNDED_DUE_TO_ISSUE, *transaction.RevocationReason, "RevocationReason")
	assertEqual(t, "abc.123", *transaction.OfferIdentifier, "OfferIdentifier")
	assertEqual(t, true, transaction.IsUpgraded, "IsUpgraded")
	assertEqual(t, OFFER_TYPE_INTRODUCTORY, transaction.OfferType, "OfferType")
	assertEqual(t, "USA", transaction.Storefront, "Storefront")
	assertEqual(t, "143441", transaction.StorefrontId, "StorefrontId")
	assertEqual(t, TRANSACTION_REASON_PURCHASE, transaction.TransactionReason, "TransactionReason")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, transaction.Environment, "Environment")
	assertEqual(t, int64(10990), transaction.Price, "Price")
	assertEqual(t, "USD", transaction.Currency, "Currency")
	assertEqual(t, OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO, transaction.OfferDiscountType, "OfferDiscountType")
}

// Test transaction with revocation
func TestVerifyAndDecodeSignedTransactionWithRevocation(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedTransaction, err := createSignedDataFromJSON("models/signedTransactionWithRevocation.json")
	assertNoError(t, err, "Failed to create signed data")

	transaction, err := verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	assertNoError(t, err, "Failed to verify and decode transaction")

	// Verify revocation fields
	assertEqual(t, REVOCATION_TYPE_REFUND_PRORATED, transaction.RevocationType, "RevocationType")
	assertEqual(t, int32(50000), transaction.RevocationPercentage, "RevocationPercentage")
}

// Test renewal info decoding
func TestVerifyAndDecodeRenewalInfo(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedRenewalInfo, err := createSignedDataFromJSON("models/signedRenewalInfo.json")
	assertNoError(t, err, "Failed to create signed data")

	renewalInfo, err := verifier.VerifyAndDecodeRenewalInfo(signedRenewalInfo)
	assertNoError(t, err, "Failed to verify and decode renewal info")

	// Verify all fields
	assertEqual(t, EXPIRATION_INTENT_CUSTOMER_CANCELLED, *renewalInfo.ExpirationIntent, "ExpirationIntent")
	assertEqual(t, "12345", renewalInfo.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, "com.example.product.2", renewalInfo.AutoRenewProductId, "AutoRenewProductId")
	assertEqual(t, "com.example.product", renewalInfo.ProductId, "ProductId")
	assertEqual(t, AUTO_RENEW_STATUS_ON, renewalInfo.AutoRenewStatus, "AutoRenewStatus")
	assertEqual(t, true, *renewalInfo.IsInBillingRetryPeriod, "IsInBillingRetryPeriod")
	assertEqual(t, PRICE_INCREASE_STATUS_CUSTOMER_HAS_NOT_RESPONDED, *renewalInfo.PriceIncreaseStatus, "PriceIncreaseStatus")
	assertEqual(t, Timestamp(1698148900000), *renewalInfo.GracePeriodExpiresDate, "GracePeriodExpiresDate")
	assertEqual(t, OFFER_TYPE_PROMOTIONAL, *renewalInfo.OfferType, "OfferType")
	assertEqual(t, "abc.123", *renewalInfo.OfferIdentifier, "OfferIdentifier")
	assertEqual(t, Timestamp(1698148800000), renewalInfo.SignedDate, "SignedDate")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, renewalInfo.Environment, "Environment")
	assertEqual(t, Timestamp(1698148800000), renewalInfo.RecentSubscriptionStartDate, "RecentSubscriptionStartDate")
	assertEqual(t, Timestamp(1698148850000), renewalInfo.RenewalDate, "RenewalDate")
	assertEqual(t, int64(9990), renewalInfo.RenewalPrice, "RenewalPrice")
	assertEqual(t, "USD", renewalInfo.Currency, "Currency")
	assertEqual(t, OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO, renewalInfo.OfferDiscountType, "OfferDiscountType")
}

// Test notification decoding (SUBSCRIBED type)
func TestVerifyAndDecodeNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	// Verify notification fields
	assertEqual(t, NOTIFICATION_TYPE_SUBSCRIBED, notification.NotificationType, "NotificationType")
	assertNotNil(t, notification.Subtype, "Subtype")
	assertEqual(t, SUBTYPE_INITIAL_BUY, *notification.Subtype, "Subtype")
	assertEqual(t, "002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assertEqual(t, "2.0", notification.Version, "Version")
	assertEqual(t, Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assertNotNil(t, notification.Data, "Data")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, notification.Data.Environment, "Data.Environment")
	assertEqual(t, int64(41234), notification.Data.AppAppleId, "Data.AppAppleId")
	assertEqual(t, "com.example", notification.Data.BundleId, "Data.BundleId")
	assertEqual(t, "1.2.3", notification.Data.BundleVersion, "Data.BundleVersion")
	assertEqual(t, "signed_transaction_info_value", notification.Data.SignedTransactionInfo, "Data.SignedTransactionInfo")
	assertEqual(t, "signed_renewal_info_value", notification.Data.SignedRenewalInfo, "Data.SignedRenewalInfo")
	assertEqual(t, STATUS_ACTIVE, notification.Data.Status, "Data.Status")

	// Verify nil fields
	assertNil(t, notification.Summary, "Summary")
	assertNil(t, notification.ExternalPurchaseToken, "ExternalPurchaseToken")
}

// Test consumption request notification
func TestVerifyAndDecodeConsumptionRequestNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedConsumptionRequestNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	assertEqual(t, NOTIFICATION_TYPE_CONSUMPTION_REQUEST, notification.NotificationType, "NotificationType")

	assertNotNil(t, notification.Data, "Data")
	assertEqual(t, CONSUMPTION_REQUEST_REASON_UNINTENDED_PURCHASE, *notification.Data.ConsumptionRequestReason, "ConsumptionRequestReason")
}

// Test summary notification
func TestVerifyAndDecodeSummaryNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedSummaryNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	assertEqual(t, NOTIFICATION_TYPE_RENEWAL_EXTENSION, notification.NotificationType, "NotificationType")
	assertNotNil(t, notification.Subtype, "Subtype")
	assertEqual(t, SUBTYPE_SUMMARY, *notification.Subtype, "Subtype")

	assertNil(t, notification.Data, "Data")
	assertNotNil(t, notification.Summary, "Summary")

	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, notification.Summary.Environment, "Summary.Environment")
	assertEqual(t, int64(41234), notification.Summary.AppAppleId, "Summary.AppAppleId")
	assertEqual(t, "com.example", notification.Summary.BundleId, "Summary.BundleId")
	assertEqual(t, "com.example.product", notification.Summary.ProductId, "Summary.ProductId")
	assertEqual(t, int64(5), notification.Summary.SucceededCount, "Summary.SucceededCount")
	assertEqual(t, int64(2), notification.Summary.FailedCount, "Summary.FailedCount")
}

func TestVerifyAndDecodeExternalPurchaseTokenNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")
	verifier.allowAnyEnvironment = true

	signedExternalPurchaseTokenNotification, err := createSignedDataFromJSON("models/signedExternalPurchaseTokenNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedExternalPurchaseTokenNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	assertEqual(t, NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION, notification.NotificationType, "NotificationType")
	assertNotNil(t, notification.Subtype, "Subtype")
	assertEqual(t, SUBTYPE_UNREPORTED, *notification.Subtype, "Subtype")
	assertEqual(t, "002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assertEqual(t, "2.0", notification.Version, "Version")
	assertEqual(t, Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assertNil(t, notification.Data, "Data")
	assertNil(t, notification.Summary, "Summary")
	assertNotNil(t, notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	assertEqual(t, "b2158121-7af9-49d4-9561-1f588205523e", notification.ExternalPurchaseToken.ExternalPurchaseId, "ExternalPurchaseToken.ExternalPurchaseId")
	assertEqual(t, Timestamp(1698148950000), notification.ExternalPurchaseToken.TokenCreationDate, "ExternalPurchaseToken.TokenCreationDate")
	assertEqual(t, int64(55555), notification.ExternalPurchaseToken.AppAppleId, "ExternalPurchaseToken.AppAppleId")
	assertEqual(t, "com.example", notification.ExternalPurchaseToken.BundleId, "ExternalPurchaseToken.BundleId")
}

func TestVerifyAndDecodeExternalPurchaseTokenSandboxNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")
	verifier.allowAnyEnvironment = true

	signedExternalPurchaseTokenNotification, err := createSignedDataFromJSON("models/signedExternalPurchaseTokenSandboxNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedExternalPurchaseTokenNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	assertEqual(t, NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION, notification.NotificationType, "NotificationType")
	assertNotNil(t, notification.Subtype, "Subtype")
	assertEqual(t, SUBTYPE_UNREPORTED, *notification.Subtype, "Subtype")
	assertEqual(t, "002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assertEqual(t, "2.0", notification.Version, "Version")
	assertEqual(t, Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	assertNil(t, notification.Data, "Data")
	assertNil(t, notification.Summary, "Summary")
	assertNotNil(t, notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	assertEqual(t, "SANDBOX_b2158121-7af9-49d4-9561-1f588205523e", notification.ExternalPurchaseToken.ExternalPurchaseId, "ExternalPurchaseToken.ExternalPurchaseId")
	assertEqual(t, Timestamp(1698148950000), notification.ExternalPurchaseToken.TokenCreationDate, "ExternalPurchaseToken.TokenCreationDate")
	assertEqual(t, int64(55555), notification.ExternalPurchaseToken.AppAppleId, "ExternalPurchaseToken.AppAppleId")
	assertEqual(t, "com.example", notification.ExternalPurchaseToken.BundleId, "ExternalPurchaseToken.BundleId")
}

// Test rescind consent notification
func TestVerifyAndDecodeRescindConsentNotification(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedNotification, err := createSignedDataFromJSON("models/signedRescindConsentNotification.json")
	assertNoError(t, err, "Failed to create signed data")

	notification, err := verifier.VerifyAndDecodeNotification(signedNotification)
	assertNoError(t, err, "Failed to verify and decode notification")

	// Verify notification type
	assertEqual(t, NOTIFICATION_TYPE_RESCIND_CONSENT, notification.NotificationType, "NotificationType")
	assertEqual(t, "002e14d5-51f5-4503-b5a8-c3a1af68eb20", notification.NotificationUUID, "NotificationUUID")
	assertEqual(t, "2.0", notification.Version, "Version")
	assertEqual(t, Timestamp(1698148900000), notification.SignedDate, "SignedDate")

	// Verify nil fields
	assertNil(t, notification.Subtype, "Subtype")
	assertNil(t, notification.Data, "Data")
	assertNil(t, notification.Summary, "Summary")
	assertNil(t, notification.ExternalPurchaseToken, "ExternalPurchaseToken")

	// Verify AppData - THIS TESTS AppData unmarshal!
	assertNotNil(t, notification.AppData, "AppData")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, notification.AppData.Environment, "AppData.Environment")
	assertEqual(t, int64(41234), notification.AppData.AppAppleId, "AppData.AppAppleId")
	assertEqual(t, "com.example", notification.AppData.BundleId, "AppData.BundleId")
	assertEqual(t, "signed_app_transaction_info_value", notification.AppData.SignedAppTransactionInfo, "AppData.SignedAppTransactionInfo")
}

// Test realtime request decoding
func TestVerifyAndDecodeRealtimeRequest(t *testing.T) {
	verifier, err := createDefaultTestSignedDataVerifier()
	assertNoError(t, err, "Failed to create verifier")

	signedRealtimeRequest, err := createSignedDataFromJSON("models/decodedRealtimeRequest.json")
	assertNoError(t, err, "Failed to create signed data")

	request, err := verifier.VerifyAndDecodeRealtimeRequest(signedRealtimeRequest)
	assertNoError(t, err, "Failed to verify and decode realtime request")

	// Verify fields - THIS TESTS DecodedRealtimeRequestBody unmarshal!
	assertEqual(t, "99371282", request.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, int64(531412), request.AppAppleId, "AppAppleId")
	assertEqual(t, "com.example.product", request.ProductId, "ProductId")
	assertEqual(t, "en-US", request.UserLocale, "UserLocale")
	assertEqual(t, "3db5c98d-8acf-4e29-831e-8e1f82f9f6e9", request.RequestIdentifier, "RequestIdentifier")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, request.Environment, "Environment")
	assertEqual(t, Timestamp(1698148900000), request.SignedDate, "SignedDate")
}
