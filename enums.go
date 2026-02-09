package appstore

// Environment is the server environment, either sandbox or production.
//
// https://developer.apple.com/documentation/appstoreserverapi/environment
type Environment string

const (
	ENVIRONMENT_PRODUCTION    Environment = "Production"
	ENVIRONMENT_SANDBOX       Environment = "Sandbox"
	ENVIRONMENT_XCODE         Environment = "Xcode"
	ENVIRONMENT_LOCAL_TESTING Environment = "LocalTesting" // Used for unit testing
)

// InAppOwnershipType is the relationship of the user with the family-shared purchase to which they have access.
//
// https://developer.apple.com/documentation/appstoreserverapi/inappownershiptype
type InAppOwnershipType string

const (
	IN_APP_OWNERSHIP_TYPE_FAMILY_SHARED InAppOwnershipType = "FAMILY_SHARED"
	IN_APP_OWNERSHIP_TYPE_PURCHASED     InAppOwnershipType = "PURCHASED"
)

// OfferType is the type of offer.
//
// https://developer.apple.com/documentation/appstoreserverapi/offertype
type OfferType int32

const (
	OFFER_TYPE_INTRODUCTORY OfferType = 1
	OFFER_TYPE_PROMOTIONAL  OfferType = 2
	OFFER_TYPE_OFFER_CODE   OfferType = 3
	OFFER_TYPE_WIN_BACK     OfferType = 4
)

// Status is the status of an auto-renewable subscription.
//
// https://developer.apple.com/documentation/appstoreserverapi/status
type Status int32

const (
	STATUS_ACTIVE               Status = 1
	STATUS_EXPIRED              Status = 2
	STATUS_BILLING_RETRY        Status = 3
	STATUS_BILLING_GRACE_PERIOD Status = 4
	STATUS_REVOKED              Status = 5
)

// Type is the type of in-app purchase products you can offer in your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/type
type Type string

const (
	TYPE_AUTO_RENEWABLE_SUBSCRIPTION Type = "Auto-Renewable Subscription"
	TYPE_NON_CONSUMABLE              Type = "Non-Consumable"
	TYPE_CONSUMABLE                  Type = "Consumable"
	TYPE_NON_RENEWING_SUBSCRIPTION   Type = "Non-Renewing Subscription"
)

// PurchasePlatform is the platform on which the customer originally purchased the app.
//
// https://developer.apple.com/documentation/storekit/apptransaction/originalplatform
type PurchasePlatform string

const (
	PURCHASE_PLATFORM_MACOS    PurchasePlatform = "macOS"
	PURCHASE_PLATFORM_IOS      PurchasePlatform = "iOS"
	PURCHASE_PLATFORM_TVOS     PurchasePlatform = "tvOS"
	PURCHASE_PLATFORM_VISIONOS PurchasePlatform = "visionOS"
)

// NotificationTypeV2 is the type that describes the in-app purchase or external purchase event for which the App Store sends the version 2 notification.
//
// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
type NotificationTypeV2 string

const (
	NOTIFICATION_TYPE_SUBSCRIBED                           NotificationTypeV2 = "SUBSCRIBED"
	NOTIFICATION_TYPE_DID_CHANGE_RENEWAL_PREF              NotificationTypeV2 = "DID_CHANGE_RENEWAL_PREF"
	NOTIFICATION_TYPE_DID_CHANGE_RENEWAL_STATUS            NotificationTypeV2 = "DID_CHANGE_RENEWAL_STATUS"
	NOTIFICATION_TYPE_OFFER_REDEEMED                       NotificationTypeV2 = "OFFER_REDEEMED"
	NOTIFICATION_TYPE_DID_RENEW                            NotificationTypeV2 = "DID_RENEW"
	NOTIFICATION_TYPE_EXPIRED                              NotificationTypeV2 = "EXPIRED"
	NOTIFICATION_TYPE_DID_FAIL_TO_RENEW                    NotificationTypeV2 = "DID_FAIL_TO_RENEW"
	NOTIFICATION_TYPE_GRACE_PERIOD_EXPIRED                 NotificationTypeV2 = "GRACE_PERIOD_EXPIRED"
	NOTIFICATION_TYPE_PRICE_INCREASE                       NotificationTypeV2 = "PRICE_INCREASE"
	NOTIFICATION_TYPE_REFUND                               NotificationTypeV2 = "REFUND"
	NOTIFICATION_TYPE_REFUND_DECLINED                      NotificationTypeV2 = "REFUND_DECLINED"
	NOTIFICATION_TYPE_CONSUMPTION_REQUEST                  NotificationTypeV2 = "CONSUMPTION_REQUEST"
	NOTIFICATION_TYPE_RENEWAL_EXTENDED                     NotificationTypeV2 = "RENEWAL_EXTENDED"
	NOTIFICATION_TYPE_REVOKE                               NotificationTypeV2 = "REVOKE"
	NOTIFICATION_TYPE_TEST                                 NotificationTypeV2 = "TEST"
	NOTIFICATION_TYPE_RENEWAL_EXTENSION                    NotificationTypeV2 = "RENEWAL_EXTENSION"
	NOTIFICATION_TYPE_REFUND_REVERSED                      NotificationTypeV2 = "REFUND_REVERSED"
	NOTIFICATION_TYPE_EXTERNAL_PURCHASE_TOKEN_NOTIFICATION NotificationTypeV2 = "EXTERNAL_PURCHASE_TOKEN"
	NOTIFICATION_TYPE_ONE_TIME_CHARGE                      NotificationTypeV2 = "ONE_TIME_CHARGE"
	NOTIFICATION_TYPE_RESCIND_CONSENT                      NotificationTypeV2 = "RESCIND_CONSENT"
)

// Subtype is a string that provides details about select notification types in version 2.
//
// https://developer.apple.com/documentation/appstoreservernotifications/subtype
type Subtype string

const (
	SUBTYPE_INITIAL_BUY          Subtype = "INITIAL_BUY"
	SUBTYPE_RESUBSCRIBE          Subtype = "RESUBSCRIBE"
	SUBTYPE_DOWNGRADE            Subtype = "DOWNGRADE"
	SUBTYPE_UPGRADE              Subtype = "UPGRADE"
	SUBTYPE_AUTO_RENEW_ENABLED   Subtype = "AUTO_RENEW_ENABLED"
	SUBTYPE_AUTO_RENEW_DISABLED  Subtype = "AUTO_RENEW_DISABLED"
	SUBTYPE_VOLUNTARY            Subtype = "VOLUNTARY"
	SUBTYPE_BILLING_RETRY        Subtype = "BILLING_RETRY"
	SUBTYPE_PRICE_INCREASE       Subtype = "PRICE_INCREASE"
	SUBTYPE_GRACE_PERIOD         Subtype = "GRACE_PERIOD"
	SUBTYPE_PENDING              Subtype = "PENDING"
	SUBTYPE_ACCEPTED             Subtype = "ACCEPTED"
	SUBTYPE_BILLING_RECOVERY     Subtype = "BILLING_RECOVERY"
	SUBTYPE_PRODUCT_NOT_FOR_SALE Subtype = "PRODUCT_NOT_FOR_SALE"
	SUBTYPE_SUMMARY              Subtype = "SUMMARY"
	SUBTYPE_FAILURE              Subtype = "FAILURE"
	SUBTYPE_UNREPORTED           Subtype = "UNREPORTED"
)

// ExpirationIntent is the reason an auto-renewable subscription expired.
//
// https://developer.apple.com/documentation/appstoreserverapi/expirationintent
type ExpirationIntent int32

const (
	EXPIRATION_INTENT_CUSTOMER_CANCELLED                         ExpirationIntent = 1
	EXPIRATION_INTENT_BILLING_ERROR                              ExpirationIntent = 2
	EXPIRATION_INTENT_CUSTOMER_DID_NOT_CONSENT_TO_PRICE_INCREASE ExpirationIntent = 3
	EXPIRATION_INTENT_PRODUCT_NOT_AVAILABLE                      ExpirationIntent = 4
	EXPIRATION_INTENT_OTHER                                      ExpirationIntent = 5
)

// ExtendReasonCode is the reason code for the subscription date extension.
//
// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
type ExtendReasonCode int32

const (
	UNDECLARED            ExtendReasonCode = 0
	CUSTOMER_SATISFACTION ExtendReasonCode = 1
	OTHER_REASON          ExtendReasonCode = 2
	SERVICE_ISSUE         ExtendReasonCode = 3
)

// OrderLookupStatus is a value that indicates whether the order ID in the request is valid for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/orderlookupstatus
type OrderLookupStatus int32

const (
	ORDER_LOOKUP_VALID   OrderLookupStatus = 0
	ORDER_LOOKUP_INVALID OrderLookupStatus = 1
)

// AutoRenewStatus is the renewal status for an auto-renewable subscription.
//
// https://developer.apple.com/documentation/appstoreserverapi/autorenewstatus
type AutoRenewStatus int32

const (
	AUTO_RENEW_STATUS_OFF AutoRenewStatus = 0
	AUTO_RENEW_STATUS_ON  AutoRenewStatus = 1
)

// PriceIncreaseStatus is the status that indicates whether an auto-renewable subscription is subject to a price increase.
//
// https://developer.apple.com/documentation/appstoreserverapi/priceincreasestatus
type PriceIncreaseStatus int32

const (
	PRICE_INCREASE_STATUS_CUSTOMER_HAS_NOT_RESPONDED                                 PriceIncreaseStatus = 0
	PRICE_INCREASE_STATUS_CUSTOMER_CONSENTED_OR_WAS_NOTIFIED_WITHOUT_NEEDING_CONSENT PriceIncreaseStatus = 1
)

// DeliveryStatus is a value that indicates whether the app successfully delivered an in-app purchase that works properly.
//
// https://developer.apple.com/documentation/appstoreserverapi/deliverystatus
type DeliveryStatus int32

const (
	DELIVERED_AND_WORKING_PROPERLY DeliveryStatus = 0
	DID_NOT_DELIVER_DUE_TO_ISSUE   DeliveryStatus = 1
	OTHER_DELIVERY_STATUS          DeliveryStatus = 2
	DELIVERED_AND_HAS_ISSUE        DeliveryStatus = 3
)

// RefundPreference is a value that indicates your preferred outcome for the refund request.
//
// https://developer.apple.com/documentation/appstoreserverapi/refundpreference
type RefundPreference int32

const (
	UNDECLARED_PREFERENCE RefundPreference = 0
	PREFER_REFUND         RefundPreference = 1
	PREFER_NO_REFUND      RefundPreference = 2
	NO_PREFERENCE         RefundPreference = 3
)

// SendAttemptResult is the success or error information the App Store server records when it attempts to send an App Store server notification to your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/sendattemptresult
type SendAttemptResult string

const (
	SUCCESS                         SendAttemptResult = "SUCCESS"
	TIMED_OUT                       SendAttemptResult = "TIMED_OUT"
	TLS_ISSUE                       SendAttemptResult = "TLS_ISSUE"
	CIRCULAR_REDIRECT               SendAttemptResult = "CIRCULAR_REDIRECT"
	NO_RESPONSE                     SendAttemptResult = "NO_RESPONSE"
	SOCKET_ISSUE                    SendAttemptResult = "SOCKET_ISSUE"
	UNSUPPORTED_CHARSET             SendAttemptResult = "UNSUPPORTED_CHARSET"
	INVALID_RESPONSE                SendAttemptResult = "INVALID_RESPONSE"
	PREMATURE_CLOSE                 SendAttemptResult = "PREMATURE_CLOSE"
	UNSUCCESSFUL_HTTP_RESPONSE_CODE SendAttemptResult = "UNSUCCESSFUL_HTTP_RESPONSE_CODE"
	OTHER                           SendAttemptResult = "OTHER"
)

// ImageState is the approval state of an image.
//
// https://developer.apple.com/documentation/retentionmessaging/imagestate
type ImageState string

const (
	IMAGE_PENDING  ImageState = "PENDING"
	IMAGE_APPROVED ImageState = "APPROVED"
	IMAGE_REJECTED ImageState = "REJECTED"
)

// MessageState is the approval state of the message.
//
// https://developer.apple.com/documentation/retentionmessaging/messagestate
type MessageState string

const (
	MESSAGE_PENDING  MessageState = "PENDING"
	MESSAGE_APPROVED MessageState = "APPROVED"
	MESSAGE_REJECTED MessageState = "REJECTED"
)

// OfferDiscountType is the payment mode for a discount offer on an In-App Purchase.
//
// https://developer.apple.com/documentation/appstoreserverapi/offerdiscounttype
type OfferDiscountType string

const (
	OFFER_DISCOUNT_TYPE_FREE_TRIAL    OfferDiscountType = "FREE_TRIAL"
	OFFER_DISCOUNT_TYPE_PAY_AS_YOU_GO OfferDiscountType = "PAY_AS_YOU_GO"
	OFFER_DISCOUNT_TYPE_PAY_UP_FRONT  OfferDiscountType = "PAY_UP_FRONT"
	OFFER_DISCOUNT_TYPE_ONE_TIME      OfferDiscountType = "ONE_TIME"
)

// RevocationReason is the reason for a refunded transaction.
//
// https://developer.apple.com/documentation/appstoreserverapi/revocationreason
type RevocationReason int32

const (
	REVOCATION_REASON_REFUNDED_FOR_OTHER_REASON RevocationReason = 0
	REVOCATION_REASON_REFUNDED_DUE_TO_ISSUE     RevocationReason = 1
)

// RevocationType is the type of the refund or revocation that applies to the transaction.
//
// https://developer.apple.com/documentation/appstoreservernotifications/revocationtype
type RevocationType string

const (
	REVOCATION_TYPE_REFUND_FULL     RevocationType = "REFUND_FULL"
	REVOCATION_TYPE_REFUND_PRORATED RevocationType = "REFUND_PRORATED"
	REVOCATION_TYPE_FAMILY_REVOKE   RevocationType = "FAMILY_REVOKE"
)

// TransactionReason is the cause of a purchase transaction, which indicates whether it’s a customer’s purchase or a renewal for an auto-renewable subscription that the system initiates.
//
// https://developer.apple.com/documentation/appstoreserverapi/transactionreason
type TransactionReason string

const (
	TRANSACTION_REASON_PURCHASE TransactionReason = "PURCHASE"
	TRANSACTION_REASON_RENEWAL  TransactionReason = "RENEWAL"
)

// ConsumptionRequestReason is the reason the customer requested the refund.
//
// https://developer.apple.com/documentation/appstoreservernotifications/consumptionrequestreason
type ConsumptionRequestReason string

const (
	CONSUMPTION_REQUEST_REASON_UNINTENDED_PURCHASE       ConsumptionRequestReason = "UNINTENDED_PURCHASE"
	CONSUMPTION_REQUEST_REASON_FULFILLMENT_ISSUE         ConsumptionRequestReason = "FULFILLMENT_ISSUE"
	CONSUMPTION_REQUEST_REASON_UNSATISFIED_WITH_PURCHASE ConsumptionRequestReason = "UNSATISFIED_WITH_PURCHASE"
	CONSUMPTION_REQUEST_REASON_LEGAL_REASON              ConsumptionRequestReason = "LEGAL_REASON"
	CONSUMPTION_REQUEST_REASON_OTHER                     ConsumptionRequestReason = "OTHER"
)

// GetTransactionHistoryVersion is the version of the Get Transaction History endpoint.
type GetTransactionHistoryVersion string

const (
	GET_TRANSACTION_HISTORY_VERSION_V1 GetTransactionHistoryVersion = "v1"
	GET_TRANSACTION_HISTORY_VERSION_V2 GetTransactionHistoryVersion = "v2"
)

// APIError is an error that indicates an invalid request.
//
// https://developer.apple.com/documentation/appstoreserverapi/error_codes
type APIError int32

const (
	API_ERROR_GENERAL_BAD_REQUEST                                 APIError = 4000000
	API_ERROR_INVALID_APP_IDENTIFIER                              APIError = 4000002
	API_ERROR_INVALID_REQUEST_REVISION                            APIError = 4000005
	API_ERROR_INVALID_TRANSACTION_ID                              APIError = 4000006
	API_ERROR_INVALID_ORIGINAL_TRANSACTION_ID                     APIError = 4000008
	API_ERROR_INVALID_EXTEND_BY_DAYS                              APIError = 4000009
	API_ERROR_INVALID_EXTEND_REASON_CODE                          APIError = 4000010
	API_ERROR_INVALID_REQUEST_IDENTIFIER                          APIError = 4000011
	API_ERROR_START_DATE_TOO_FAR_IN_PAST                          APIError = 4000012
	API_ERROR_START_DATE_AFTER_END_DATE                           APIError = 4000013
	API_ERROR_INVALID_PAGINATION_TOKEN                            APIError = 4000014
	API_ERROR_INVALID_START_DATE                                  APIError = 4000015
	API_ERROR_INVALID_END_DATE                                    APIError = 4000016
	API_ERROR_PAGINATION_TOKEN_EXPIRED                            APIError = 4000017
	API_ERROR_INVALID_NOTIFICATION_TYPE                           APIError = 4000018
	API_ERROR_MULTIPLE_FILTERS_SUPPLIED                           APIError = 4000019
	API_ERROR_INVALID_TEST_NOTIFICATION_TOKEN                     APIError = 4000020
	API_ERROR_INVALID_SORT                                        APIError = 4000021
	API_ERROR_INVALID_PRODUCT_TYPE                                APIError = 4000022
	API_ERROR_INVALID_PRODUCT_ID                                  APIError = 4000023
	API_ERROR_INVALID_SUBSCRIPTION_GROUP_IDENTIFIER               APIError = 4000024
	API_ERROR_INVALID_EXCLUDE_REVOKED                             APIError = 4000025
	API_ERROR_INVALID_IN_APP_OWNERSHIP_TYPE                       APIError = 4000026
	API_ERROR_INVALID_EMPTY_STOREFRONT_COUNTRY_CODE_LIST          APIError = 4000027
	API_ERROR_INVALID_STOREFRONT_COUNTRY_CODE                     APIError = 4000028
	API_ERROR_INVALID_REVOKED                                     APIError = 4000030
	API_ERROR_INVALID_STATUS                                      APIError = 4000031
	API_ERROR_INVALID_ACCOUNT_TENURE                              APIError = 4000032
	API_ERROR_INVALID_APP_ACCOUNT_TOKEN                           APIError = 4000033
	API_ERROR_INVALID_CONSUMPTION_STATUS                          APIError = 4000034
	API_ERROR_INVALID_CUSTOMER_CONSENTED                          APIError = 4000035
	API_ERROR_INVALID_DELIVERY_STATUS                             APIError = 4000036
	API_ERROR_INVALID_LIFETIME_DOLLARS_PURCHASED                  APIError = 4000037
	API_ERROR_INVALID_LIFETIME_DOLLARS_REFUNDED                   APIError = 4000038
	API_ERROR_INVALID_PLATFORM                                    APIError = 4000039
	API_ERROR_INVALID_PLAY_TIME                                   APIError = 4000040
	API_ERROR_INVALID_SAMPLE_CONTENT_PROVIDED                     APIError = 4000041
	API_ERROR_INVALID_USER_STATUS                                 APIError = 4000042
	API_ERROR_INVALID_TRANSACTION_NOT_CONSUMABLE                  APIError = 4000043
	API_ERROR_INVALID_TRANSACTION_TYPE_NOT_SUPPORTED              APIError = 4000047
	API_ERROR_APP_TRANSACTION_ID_NOT_SUPPORTED_ERROR              APIError = 4000048
	API_ERROR_INVALID_IMAGE                                       APIError = 4000161
	API_ERROR_HEADER_TOO_LONG                                     APIError = 4000162
	API_ERROR_BODY_TOO_LONG                                       APIError = 4000163
	API_ERROR_INVALID_LOCALE                                      APIError = 4000164
	API_ERROR_ALT_TEXT_TOO_LONG                                   APIError = 4000175
	API_ERROR_INVALID_APP_ACCOUNT_TOKEN_UUID_ERROR                APIError = 4000183
	API_ERROR_FAMILY_TRANSACTION_NOT_SUPPORTED_ERROR              APIError = 4000185
	API_ERROR_TRANSACTION_ID_IS_NOT_ORIGINAL_TRANSACTION_ID_ERROR APIError = 4000187
	API_ERROR_SUBSCRIPTION_EXTENSION_INELIGIBLE                   APIError = 4030004
	API_ERROR_SUBSCRIPTION_MAX_EXTENSION                          APIError = 4030005
	API_ERROR_FAMILY_SHARED_SUBSCRIPTION_EXTENSION_INELIGIBLE     APIError = 4030007
	API_ERROR_MAXIMUM_NUMBER_OF_IMAGES_REACHED                    APIError = 4030014
	API_ERROR_MAXIMUM_NUMBER_OF_MESSAGES_REACHED                  APIError = 4030016
	API_ERROR_MESSAGE_NOT_APPROVED                                APIError = 4030017
	API_ERROR_IMAGE_NOT_APPROVED                                  APIError = 4030018
	API_ERROR_IMAGE_IN_USE                                        APIError = 4030019
	API_ERROR_ACCOUNT_NOT_FOUND                                   APIError = 4040001
	API_ERROR_ACCOUNT_NOT_FOUND_RETRYABLE                         APIError = 4040002
	API_ERROR_APP_NOT_FOUND                                       APIError = 4040003
	API_ERROR_APP_NOT_FOUND_RETRYABLE                             APIError = 4040004
	API_ERROR_ORIGINAL_TRANSACTION_ID_NOT_FOUND                   APIError = 4040005
	API_ERROR_ORIGINAL_TRANSACTION_ID_NOT_FOUND_RETRYABLE         APIError = 4040006
	API_ERROR_SERVER_NOTIFICATION_URL_NOT_FOUND                   APIError = 4040007
	API_ERROR_TEST_NOTIFICATION_NOT_FOUND                         APIError = 4040008
	API_ERROR_STATUS_REQUEST_NOT_FOUND                            APIError = 4040009
	API_ERROR_TRANSACTION_ID_NOT_FOUND                            APIError = 4040010
	API_ERROR_IMAGE_NOT_FOUND                                     APIError = 4040014
	API_ERROR_MESSAGE_NOT_FOUND                                   APIError = 4040015
	API_ERROR_APP_TRANSACTION_DOES_NOT_EXIST_ERROR                APIError = 4040019
	API_ERROR_IMAGE_ALREADY_EXISTS                                APIError = 4090000
	API_ERROR_MESSAGE_ALREADY_EXISTS                              APIError = 4090001
	API_ERROR_RATE_LIMIT_EXCEEDED                                 APIError = 4290000
	API_ERROR_GENERAL_INTERNAL                                    APIError = 5000000
	API_ERROR_GENERAL_INTERNAL_RETRYABLE                          APIError = 5000001
)
