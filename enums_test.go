package appstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnumsExpanded(t *testing.T) {
	assert := assert.New(t)
	// Environment
	environments := []Environment{ENVIRONMENT_PRODUCTION, ENVIRONMENT_SANDBOX, ENVIRONMENT_XCODE, ENVIRONMENT_LOCAL_TESTING}
	for _, e := range environments {
		assert.Equal(true, e.IsValid(), "Environment.IsValid")
		assert.Equal(string(e), e.Raw(), "Environment.Raw")
	}
	assert.Equal(false, Environment("Invalid").IsValid(), "Environment(Invalid).IsValid")

	// InAppOwnershipType
	ownershipTypes := []InAppOwnershipType{IN_APP_OWNERSHIP_TYPE_FAMILY_SHARED, IN_APP_OWNERSHIP_TYPE_PURCHASED}
	for _, i := range ownershipTypes {
		assert.Equal(true, i.IsValid(), "InAppOwnershipType.IsValid")
		assert.Equal(string(i), i.Raw(), "InAppOwnershipType.Raw")
	}
	assert.Equal(false, InAppOwnershipType("Invalid").IsValid(), "InAppOwnershipType(Invalid).IsValid")

	// OfferType
	offerTypes := []OfferType{OFFER_TYPE_INTRODUCTORY, OFFER_TYPE_PROMOTIONAL, OFFER_TYPE_OFFER_CODE, OFFER_TYPE_WIN_BACK}
	for _, o := range offerTypes {
		assert.Equal(true, o.IsValid(), "OfferType.IsValid")
		assert.Equal(int32(o), o.Raw(), "OfferType.Raw")
	}
	assert.Equal(false, OfferType(99).IsValid(), "OfferType(99).IsValid")

	// Status
	statuses := []Status{STATUS_ACTIVE, STATUS_EXPIRED, STATUS_BILLING_RETRY, STATUS_BILLING_GRACE_PERIOD, STATUS_REVOKED}
	for _, s := range statuses {
		assert.Equal(true, s.IsValid(), "Status.IsValid")
		assert.Equal(int32(s), s.Raw(), "Status.Raw")
	}
	assert.Equal(false, Status(99).IsValid(), "Status(99).IsValid")

	// Type
	types := []Type{TYPE_AUTO_RENEWABLE_SUBSCRIPTION, TYPE_NON_CONSUMABLE, TYPE_CONSUMABLE, TYPE_NON_RENEWING_SUBSCRIPTION}
	for _, typ := range types {
		assert.Equal(true, typ.IsValid(), "Type.IsValid")
		assert.Equal(string(typ), typ.Raw(), "Type.Raw")
	}
	assert.Equal(false, Type("Invalid").IsValid(), "Type(Invalid).IsValid")

	// PurchasePlatform
	platforms := []PurchasePlatform{PURCHASE_PLATFORM_MACOS, PURCHASE_PLATFORM_IOS, PURCHASE_PLATFORM_TVOS, PURCHASE_PLATFORM_VISIONOS}
	for _, p := range platforms {
		assert.Equal(true, p.IsValid(), "PurchasePlatform.IsValid")
		assert.Equal(string(p), p.Raw(), "PurchasePlatform.Raw")
	}
	assert.Equal(false, PurchasePlatform("Invalid").IsValid(), "PurchasePlatform(Invalid).IsValid")

	// NotificationTypeV2
	notificationTypes := []NotificationTypeV2{
		NOTIFICATION_TYPE_SUBSCRIBED,
		NOTIFICATION_TYPE_DID_CHANGE_RENEWAL_PREF,
		NOTIFICATION_TYPE_DID_CHANGE_RENEWAL_STATUS,
		NOTIFICATION_TYPE_OFFER_REDEEMED,
		NOTIFICATION_TYPE_DID_RENEW,
		NOTIFICATION_TYPE_EXPIRED,
		NOTIFICATION_TYPE_DID_FAIL_TO_RENEW,
		NOTIFICATION_TYPE_GRACE_PERIOD_EXPIRED,
		NOTIFICATION_TYPE_PRICE_INCREASE,
		NOTIFICATION_TYPE_REFUND,
		NOTIFICATION_TYPE_REFUND_DECLINED,
		NOTIFICATION_TYPE_CONSUMPTION_REQUEST,
		NOTIFICATION_TYPE_RENEWAL_EXTENDED,
		NOTIFICATION_TYPE_REVOKE,
		NOTIFICATION_TYPE_TEST,
		NOTIFICATION_TYPE_RENEWAL_EXTENSION,
		NOTIFICATION_TYPE_REFUND_REVERSED,
		NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION,
		NOTIFICATION_TYPE_ONE_TIME_CHARGE,
		NOTIFICATION_TYPE_RESCIND_CONSENT,
	}
	for _, n := range notificationTypes {
		assert.Equal(true, n.IsValid(), "NotificationTypeV2.IsValid")
		assert.Equal(string(n), n.Raw(), "NotificationTypeV2.Raw")
	}
	assert.Equal(false, NotificationTypeV2("Invalid").IsValid(), "NotificationTypeV2(Invalid).IsValid")

	// Subtype
	subtypes := []Subtype{
		SUBTYPE_INITIAL_BUY,
		SUBTYPE_RESUBSCRIBE,
		SUBTYPE_DOWNGRADE,
		SUBTYPE_UPGRADE,
		SUBTYPE_AUTO_RENEW_ENABLED,
		SUBTYPE_AUTO_RENEW_DISABLED,
		SUBTYPE_VOLUNTARY,
		SUBTYPE_BILLING_RETRY,
		SUBTYPE_PRICE_INCREASE,
		SUBTYPE_GRACE_PERIOD,
		SUBTYPE_PENDING,
		SUBTYPE_ACCEPTED,
		SUBTYPE_BILLING_RECOVERY,
		SUBTYPE_PRODUCT_NOT_FOR_SALE,
		SUBTYPE_SUMMARY,
		SUBTYPE_FAILURE,
		SUBTYPE_UNREPORTED,
	}
	for _, s := range subtypes {
		assert.Equal(true, s.IsValid(), "Subtype.IsValid")
		assert.Equal(string(s), s.Raw(), "Subtype.Raw")
	}
	assert.Equal(false, Subtype("Invalid").IsValid(), "Subtype(Invalid).IsValid")

	// ExpirationIntent
	expirationIntents := []ExpirationIntent{
		EXPIRATION_INTENT_CUSTOMER_CANCELLED,
		EXPIRATION_INTENT_BILLING_ERROR,
		EXPIRATION_INTENT_CUSTOMER_DID_NOT_CONSENT_TO_PRICE_INCREASE,
		EXPIRATION_INTENT_PRODUCT_NOT_AVAILABLE,
		EXPIRATION_INTENT_OTHER,
	}
	for _, e := range expirationIntents {
		assert.Equal(true, e.IsValid(), "ExpirationIntent.IsValid")
		assert.Equal(int32(e), e.Raw(), "ExpirationIntent.Raw")
	}
	assert.Equal(false, ExpirationIntent(99).IsValid(), "ExpirationIntent(99).IsValid")

	// ExtendReasonCode
	extendReasonCodes := []ExtendReasonCode{EXTEND_REASON_CODE_UNDECLARED, EXTEND_REASON_CODE_CUSTOMER_SATISFACTION, EXTEND_REASON_CODE_OTHER_REASON, EXTEND_REASON_CODE_SERVICE_ISSUE}
	for _, e := range extendReasonCodes {
		assert.Equal(true, e.IsValid(), "ExtendReasonCode.IsValid")
		assert.Equal(int32(e), e.Raw(), "ExtendReasonCode.Raw")
	}
	assert.Equal(false, ExtendReasonCode(99).IsValid(), "ExtendReasonCode(99).IsValid")

	// OrderLookupStatus
	orderLookupStatuses := []OrderLookupStatus{ORDER_LOOKUP_VALID, ORDER_LOOKUP_INVALID}
	for _, o := range orderLookupStatuses {
		assert.Equal(true, o.IsValid(), "OrderLookupStatus.IsValid")
		assert.Equal(int32(o), o.Raw(), "OrderLookupStatus.Raw")
	}
	assert.Equal(false, OrderLookupStatus(99).IsValid(), "OrderLookupStatus(99).IsValid")

	// AutoRenewStatus
	autoRenewStatuses := []AutoRenewStatus{AUTO_RENEW_STATUS_OFF, AUTO_RENEW_STATUS_ON}
	for _, a := range autoRenewStatuses {
		assert.Equal(true, a.IsValid(), "AutoRenewStatus.IsValid")
		assert.Equal(int32(a), a.Raw(), "AutoRenewStatus.Raw")
	}
	assert.Equal(false, AutoRenewStatus(99).IsValid(), "AutoRenewStatus(99).IsValid")

	// PriceIncreaseStatus
	priceIncreaseStatuses := []PriceIncreaseStatus{
		PRICE_INCREASE_STATUS_CUSTOMER_HAS_NOT_RESPONDED,
		PRICE_INCREASE_STATUS_CUSTOMER_CONSENTED_OR_WAS_NOTIFIED_WITHOUT_NEEDING_CONSENT,
	}
	for _, p := range priceIncreaseStatuses {
		assert.Equal(true, p.IsValid(), "PriceIncreaseStatus.IsValid")
		assert.Equal(int32(p), p.Raw(), "PriceIncreaseStatus.Raw")
	}
	assert.Equal(false, PriceIncreaseStatus(99).IsValid(), "PriceIncreaseStatus(99).IsValid")

	// DeliveryStatus
	deliveryStatuses := []DeliveryStatus{DELIVERY_STATUS_DELIVERED_AND_WORKING_PROPERLY, DELIVERY_STATUS_DID_NOT_DELIVER_DUE_TO_ISSUE, DELIVERY_STATUS_OTHER_DELIVERY_STATUS, DELIVERY_STATUS_DELIVERED_AND_HAS_ISSUE}
	for _, d := range deliveryStatuses {
		assert.Equal(true, d.IsValid(), "DeliveryStatus.IsValid")
		assert.Equal(int32(d), d.Raw(), "DeliveryStatus.Raw")
	}
	assert.Equal(false, DeliveryStatus(99).IsValid(), "DeliveryStatus(99).IsValid")

	// RefundPreference
	refundPreferences := []RefundPreference{REFUND_PREFERENCE_UNDECLARED, REFUND_PREFERENCE_PREFER_REFUND, REFUND_PREFERENCE_PREFER_NO_REFUND, REFUND_PREFERENCE_NO_PREFERENCE}
	for _, r := range refundPreferences {
		assert.Equal(true, r.IsValid(), "RefundPreference.IsValid")
		assert.Equal(int32(r), r.Raw(), "RefundPreference.Raw")
	}
	assert.Equal(false, RefundPreference(99).IsValid(), "RefundPreference(99).IsValid")

	// SendAttemptResult
	sendAttemptResults := []SendAttemptResult{SEND_ATTEMPT_RESULT_SUCCESS, SEND_ATTEMPT_RESULT_TIMED_OUT, SEND_ATTEMPT_RESULT_TLS_ISSUE, SEND_ATTEMPT_RESULT_CIRCULAR_REDIRECT, SEND_ATTEMPT_RESULT_NO_RESPONSE, SEND_ATTEMPT_RESULT_SOCKET_ISSUE, SEND_ATTEMPT_RESULT_UNSUPPORTED_CHARSET, SEND_ATTEMPT_RESULT_INVALID_RESPONSE, SEND_ATTEMPT_RESULT_PREMATURE_CLOSE, SEND_ATTEMPT_RESULT_UNSUCCESSFUL_HTTP_RESPONSE_CODE, SEND_ATTEMPT_RESULT_OTHER}
	for _, s := range sendAttemptResults {
		assert.Equal(true, s.IsValid(), "SendAttemptResult.IsValid")
		assert.Equal(string(s), s.Raw(), "SendAttemptResult.Raw")
	}
	assert.Equal(false, SendAttemptResult("Invalid").IsValid(), "SendAttemptResult(Invalid).IsValid")

	// ImageState
	imageStates := []ImageState{IMAGE_STATE_PENDING, IMAGE_STATE_APPROVED, IMAGE_STATE_REJECTED}
	for _, i := range imageStates {
		assert.Equal(true, i.IsValid(), "ImageState.IsValid")
		assert.Equal(string(i), i.Raw(), "ImageState.Raw")
	}
	assert.Equal(false, ImageState("Invalid").IsValid(), "ImageState(Invalid).IsValid")

	// MessageState
	messageStates := []MessageState{MESSAGE_STATE_PENDING, MESSAGE_STATE_APPROVED, MESSAGE_STATE_REJECTED}
	for _, m := range messageStates {
		assert.Equal(true, m.IsValid(), "MessageState.IsValid")
		assert.Equal(string(m), m.Raw(), "MessageState.Raw")
	}
	assert.Equal(false, MessageState("Invalid").IsValid(), "MessageState(Invalid).IsValid")

	// OfferDiscountType
	offerDiscountTypes := []OfferDiscountType{OFFER_DISCOUNT_TYPE_FREE_TRIAL, OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO, OFFER_DISCOUNT_TYPE_PAY_UP_FRONT, OFFER_DISCOUNT_TYPE_ONE_TIME}
	for _, o := range offerDiscountTypes {
		assert.Equal(true, o.IsValid(), "OfferDiscountType.IsValid")
		assert.Equal(string(o), o.Raw(), "OfferDiscountType.Raw")
	}
	assert.Equal(false, OfferDiscountType("Invalid").IsValid(), "OfferDiscountType(Invalid).IsValid")

	// RevocationReason
	revocationReasons := []RevocationReason{REVOCATION_REASON_REFUNDED_FOR_OTHER_REASON, REVOCATION_REASON_REFUNDED_DUE_TO_ISSUE}
	for _, r := range revocationReasons {
		assert.Equal(true, r.IsValid(), "RevocationReason.IsValid")
		assert.Equal(int32(r), r.Raw(), "RevocationReason.Raw")
	}
	assert.Equal(false, RevocationReason(99).IsValid(), "RevocationReason(99).IsValid")

	// RevocationType
	revocationTypes := []RevocationType{REVOCATION_TYPE_REFUND_FULL, REVOCATION_TYPE_REFUND_PRORATED, REVOCATION_TYPE_FAMILY_REVOKE}
	for _, r := range revocationTypes {
		assert.Equal(true, r.IsValid(), "RevocationType.IsValid")
		assert.Equal(string(r), r.Raw(), "RevocationType.Raw")
	}
	assert.Equal(false, RevocationType("Invalid").IsValid(), "RevocationType(Invalid).IsValid")

	// TransactionReason
	transactionReasons := []TransactionReason{TRANSACTION_REASON_PURCHASE, TRANSACTION_REASON_RENEWAL}
	for _, t_reason := range transactionReasons {
		assert.Equal(true, t_reason.IsValid(), "TransactionReason.IsValid")
		assert.Equal(string(t_reason), t_reason.Raw(), "TransactionReason.Raw")
	}
	assert.Equal(false, TransactionReason("Invalid").IsValid(), "TransactionReason(Invalid).IsValid")

	// ConsumptionRequestReason
	consumptionRequestReasons := []ConsumptionRequestReason{CONSUMPTION_REQUEST_REASON_UNINTENDED_PURCHASE, CONSUMPTION_REQUEST_REASON_FULFILLMENT_ISSUE, CONSUMPTION_REQUEST_REASON_UNSATISFIED_WITH_PURCHASE, CONSUMPTION_REQUEST_REASON_LEGAL_REASON, CONSUMPTION_REQUEST_REASON_OTHER}
	for _, c := range consumptionRequestReasons {
		assert.Equal(true, c.IsValid(), "ConsumptionRequestReason.IsValid")
		assert.Equal(string(c), c.Raw(), "ConsumptionRequestReason.Raw")
	}
	assert.Equal(false, ConsumptionRequestReason("Invalid").IsValid(), "ConsumptionRequestReason(Invalid).IsValid")

	// GetTransactionHistoryVersion
	historyVersions := []GetTransactionHistoryVersion{GET_TRANSACTION_HISTORY_VERSION_V1, GET_TRANSACTION_HISTORY_VERSION_V2}
	for _, g := range historyVersions {
		assert.Equal(true, g.IsValid(), "GetTransactionHistoryVersion.IsValid")
	}
	assert.Equal(false, GetTransactionHistoryVersion("Invalid").IsValid(), "GetTransactionHistoryVersion(Invalid).IsValid")

	// APIError
	apiErrors := []APIError{
		API_ERROR_GENERAL_BAD_REQUEST,
		API_ERROR_INVALID_APP_IDENTIFIER,
		API_ERROR_INVALID_REQUEST_REVISION,
		API_ERROR_INVALID_TRANSACTION_ID,
		API_ERROR_INVALID_ORIGINAL_TRANSACTION_ID,
		API_ERROR_INVALID_EXTEND_BY_DAYS,
		API_ERROR_INVALID_EXTEND_REASON_CODE,
		API_ERROR_INVALID_REQUEST_IDENTIFIER,
		API_ERROR_START_DATE_TOO_FAR_IN_PAST,
		API_ERROR_START_DATE_AFTER_END_DATE,
		API_ERROR_INVALID_PAGINATION_TOKEN,
		API_ERROR_INVALID_START_DATE,
		API_ERROR_INVALID_END_DATE,
		API_ERROR_PAGINATION_TOKEN_EXPIRED,
		API_ERROR_INVALID_NOTIFICATION_TYPE,
		API_ERROR_MULTIPLE_FILTERS_SUPPLIED,
		API_ERROR_INVALID_TEST_NOTIFICATION_TOKEN,
		API_ERROR_INVALID_SORT,
		API_ERROR_INVALID_PRODUCT_TYPE,
		API_ERROR_INVALID_PRODUCT_ID,
		API_ERROR_INVALID_SUBSCRIPTION_GROUP_IDENTIFIER,
		API_ERROR_INVALID_EXCLUDE_REVOKED,
		API_ERROR_INVALID_IN_APP_OWNERSHIP_TYPE,
		API_ERROR_INVALID_EMPTY_STOREFRONT_COUNTRY_CODE_LIST,
		API_ERROR_INVALID_STOREFRONT_COUNTRY_CODE,
		API_ERROR_INVALID_REVOKED,
		API_ERROR_INVALID_STATUS,
		API_ERROR_INVALID_ACCOUNT_TENURE,
		API_ERROR_INVALID_APP_ACCOUNT_TOKEN,
		API_ERROR_INVALID_CONSUMPTION_STATUS,
		API_ERROR_INVALID_CUSTOMER_CONSENTED,
		API_ERROR_INVALID_DELIVERY_STATUS,
		API_ERROR_INVALID_LIFETIME_DOLLARS_PURCHASED,
		API_ERROR_INVALID_LIFETIME_DOLLARS_REFUNDED,
		API_ERROR_INVALID_PLATFORM,
		API_ERROR_INVALID_PLAY_TIME,
		API_ERROR_INVALID_SAMPLE_CONTENT_PROVIDED,
		API_ERROR_INVALID_USER_STATUS,
		API_ERROR_INVALID_TRANSACTION_NOT_CONSUMABLE,
		API_ERROR_INVALID_TRANSACTION_TYPE_NOT_SUPPORTED,
		API_ERROR_APP_TRANSACTION_ID_NOT_SUPPORTED_ERROR,
		API_ERROR_INVALID_IMAGE,
		API_ERROR_HEADER_TOO_LONG,
		API_ERROR_BODY_TOO_LONG,
		API_ERROR_INVALID_LOCALE,
		API_ERROR_ALT_TEXT_TOO_LONG,
		API_ERROR_INVALID_APP_ACCOUNT_TOKEN_UUID_ERROR,
		API_ERROR_FAMILY_TRANSACTION_NOT_SUPPORTED_ERROR,
		API_ERROR_TRANSACTION_ID_IS_NOT_ORIGINAL_TRANSACTION_ID_ERROR,
		API_ERROR_SUBSCRIPTION_EXTENSION_INELIGIBLE,
		API_ERROR_SUBSCRIPTION_MAX_EXTENSION,
		API_ERROR_FAMILY_SHARED_SUBSCRIPTION_EXTENSION_INELIGIBLE,
		API_ERROR_MAXIMUM_NUMBER_OF_IMAGES_REACHED,
		API_ERROR_MAXIMUM_NUMBER_OF_MESSAGES_REACHED,
		API_ERROR_MESSAGE_NOT_APPROVED,
		API_ERROR_IMAGE_NOT_APPROVED,
		API_ERROR_IMAGE_IN_USE,
		API_ERROR_ACCOUNT_NOT_FOUND,
		API_ERROR_ACCOUNT_NOT_FOUND_RETRYABLE,
		API_ERROR_APP_NOT_FOUND,
		API_ERROR_APP_NOT_FOUND_RETRYABLE,
		API_ERROR_ORIGINAL_TRANSACTION_ID_NOT_FOUND,
		API_ERROR_ORIGINAL_TRANSACTION_ID_NOT_FOUND_RETRYABLE,
		API_ERROR_SERVER_NOTIFICATION_URL_NOT_FOUND,
		API_ERROR_TEST_NOTIFICATION_NOT_FOUND,
		API_ERROR_STATUS_REQUEST_NOT_FOUND,
		API_ERROR_TRANSACTION_ID_NOT_FOUND,
		API_ERROR_IMAGE_NOT_FOUND,
		API_ERROR_MESSAGE_NOT_FOUND,
		API_ERROR_APP_TRANSACTION_DOES_NOT_EXIST_ERROR,
		API_ERROR_IMAGE_ALREADY_EXISTS,
		API_ERROR_MESSAGE_ALREADY_EXISTS,
		API_ERROR_RATE_LIMIT_EXCEEDED,
		API_ERROR_GENERAL_INTERNAL,
		API_ERROR_GENERAL_INTERNAL_RETRYABLE,
	}
	for _, apiErr := range apiErrors {
		assert.Equal(true, apiErr.IsValid(), "APIError.IsValid")
		assert.Equal(int32(apiErr), apiErr.Raw(), "APIError.Raw")
	}
	assert.Equal(false, APIError(99).IsValid(), "APIError(99).IsValid")
}
