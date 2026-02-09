package appstore

import (
	"encoding/json"

	"github.com/laishere/app-store-server-library-go/internal"
)

// HistoryResponse is a response that contains the customer's transaction history for an app.
//
// https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	// A token you use in a query to request the next set of transactions for the customer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/revision
	Revision string `json:"revision,omitempty"`

	// A Boolean value indicating whether the App Store has more transaction data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/hasmore
	HasMore bool `json:"hasMore,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The unique identifier of an app in the App Store.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The server environment in which you're making the request, whether sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// An array of in-app purchase transactions for the customer, signed by Apple, in JSON Web Signature format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactions []string `json:"signedTransactions,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *HistoryResponse) UnmarshalJSON(data []byte) error {
	type Alias HistoryResponse
	aux := &struct {
		Environment any `json:"environment"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Environment, &h.Environment, &h.RawEnvironment, h.Environment.Values())

	return nil
}

// TransactionInfoResponse is a response that contains signed transaction information for a single transaction.
//
// https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type TransactionInfoResponse struct {
	// A customer’s in-app purchase transaction, signed by Apple, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactionInfo string `json:"signedTransactionInfo,omitempty"`
}

// StatusResponse is a response that contains status information for all of a customer's auto-renewable subscriptions in your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/statusresponse
type StatusResponse struct {
	// The server environment, sandbox or production, in which the App Store generated the response.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The unique identifier of an app in the App Store.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// An array of information for auto-renewable subscriptions, including App Store-signed transaction information and App Store-signed renewal information.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifieritem
	Data []SubscriptionGroupIdentifierItem `json:"data,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *StatusResponse) UnmarshalJSON(data []byte) error {
	type Alias StatusResponse
	aux := &struct {
		Environment any `json:"environment"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Environment, &s.Environment, &s.RawEnvironment, s.Environment.Values())

	return nil
}

// SubscriptionGroupIdentifierItem is information for auto-renewable subscriptions, including signed transaction information and signed renewal information, for one subscription group.
//
// https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifieritem
type SubscriptionGroupIdentifierItem struct {
	// The identifier of the subscription group that the subscription belongs to.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifier
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier,omitempty"`

	// An array of the most recent App Store-signed transaction information and App Store-signed renewal information for all auto-renewable subscriptions in the subscription group.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/lasttransactionsitem
	LastTransactions []LastTransactionsItem `json:"lastTransactions,omitempty"`
}

// LastTransactionsItem is the most recent App Store-signed transaction information and App Store-signed renewal information for an auto-renewable subscription.
//
// https://developer.apple.com/documentation/appstoreserverapi/lasttransactionsitem
type LastTransactionsItem struct {
	// The status of the auto-renewable subscription.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/status
	Status Status `json:"status,omitempty"`

	// See status
	RawStatus int32 `json:"rawStatus,omitempty"`

	// The original transaction identifier of a purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/originaltransactionid
	OriginalTransactionId string `json:"originalTransactionId,omitempty"`

	// Transaction information signed by the App Store, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactionInfo string `json:"signedTransactionInfo,omitempty"`

	// Subscription renewal information, signed by the App Store, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfo
	SignedRenewalInfo string `json:"signedRenewalInfo,omitempty"`
}

// NotificationHistoryRequest is the request body for notification history.
//
// https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	// The start date of the timespan for the requested App Store Server Notification history records. The startDate needs to precede the endDate. Choose a startDate that's within the past 180 days from the current date.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/startdate
	StartDate int64 `json:"startDate,omitempty"`

	// The end date of the timespan for the requested App Store Server Notification history records. Choose an endDate that's later than the startDate. If you choose an endDate in the future, the endpoint automatically uses the current date as the endDate.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/enddate
	EndDate int64 `json:"endDate,omitempty"`

	// A notification type. Provide this field to limit the notification history records to those with this one notification type. For a list of notifications types, see notificationType.
	// Include either the transactionId or the notificationType in your query, but not both.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/notificationtype
	NotificationType NotificationTypeV2 `json:"notificationType,omitempty"`

	// A notification subtype. Provide this field to limit the notification history records to those with this one notification subtype. For a list of subtypes, see subtype. If you specify a notificationSubtype, you need to also specify its related notificationType.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/notificationsubtype
	NotificationSubtype Subtype `json:"notificationSubtype,omitempty"`

	// The transaction identifier, which may be an original transaction identifier, of any transaction belonging to the customer. Provide this field to limit the notification history request to this one customer.
	// Include either the transactionId or the notificationType in your query, but not both.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/transactionid
	TransactionId string `json:"transactionId,omitempty"`

	// A Boolean value you set to true to request only the notifications that haven’t reached your server successfully. The response also includes notifications that the App Store server is currently retrying to send to your server.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/onlyfailures
	OnlyFailures bool `json:"onlyFailures,omitempty"`
}

// NotificationHistoryResponse is a response that contains the App Store Server Notifications history for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponse struct {
	// A pagination token that you return to the endpoint on a subsequent call to receive the next set of results.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/paginationtoken
	PaginationToken string `json:"paginationToken,omitempty"`

	// A Boolean value indicating whether the App Store has more transaction data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/hasmore
	HasMore bool `json:"hasMore,omitempty"`

	// An array of App Store server notification history records.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory,omitempty"`
}

// NotificationHistoryResponseItem is the App Store server notification history record, including the signed notification payload and the result of the server's first send attempt.
//
// https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
type NotificationHistoryResponseItem struct {
	// A cryptographically signed payload, in JSON Web Signature (JWS) format, containing the response body for a version 2 notification.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/signedpayload
	SignedPayload string `json:"signedPayload,omitempty"`

	// An array of information the App Store server records for its attempts to send a notification to your server. The maximum number of entries in the array is six.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
	SendAttempts []SendAttemptItem `json:"sendAttempts,omitempty"`
}

// SendAttemptItem is the success or error information and the date the App Store server records when it attempts to send a server notification to your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
type SendAttemptItem struct {
	// The date the App Store server attempts to send a notification.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/attemptdate
	AttemptDate int64 `json:"attemptDate,omitempty"`

	// The success or error information the App Store server records when it attempts to send an App Store server notification to your server.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/sendattemptresult
	SendAttemptResult SendAttemptResult `json:"sendAttemptResult,omitempty"`

	// See sendAttemptResult
	RawSendAttemptResult string `json:"rawSendAttemptResult,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SendAttemptItem) UnmarshalJSON(data []byte) error {
	type Alias SendAttemptItem
	aux := &struct {
		SendAttemptResult any `json:"sendAttemptResult"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.SendAttemptResult, &s.SendAttemptResult, &s.RawSendAttemptResult, s.SendAttemptResult.Values())

	return nil
}

// OrderLookupResponse is a response that includes the order lookup status and an array of signed transactions for the in-app purchases in the order.
//
// https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	// The status that indicates whether the order ID is valid.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/orderlookupstatus
	Status OrderLookupStatus `json:"status,omitempty"`

	// See status
	RawStatus int32 `json:"rawStatus,omitempty"`

	// An array of in-app purchase transactions that are part of order, signed by Apple, in JSON Web Signature format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactions []string `json:"signedTransactions,omitempty"`
}

// RefundHistoryResponse is a response that contains an array of signed JSON Web Signature (JWS) refunded transactions, and paging information.
//
// https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundHistoryResponse struct {
	// A list of up to 20 JWS transactions, or an empty array if the customer hasn't received any refunds in your app. The transactions are sorted in ascending order by revocationDate.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactions []string `json:"signedTransactions,omitempty"`

	// A token you use in a query to request the next set of transactions for the customer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/revision
	Revision string `json:"revision,omitempty"`

	// A Boolean value indicating whether the App Store has more transaction data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/hasmore
	HasMore bool `json:"hasMore,omitempty"`
}

// ExtendRenewalDateRequest is the request body that contains subscription-renewal-extension data for an individual subscription.
//
// https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldaterequest
type ExtendRenewalDateRequest struct {
	// The number of days to extend the subscription renewal date.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/extendbydays
	ExtendByDays int32 `json:"extendByDays,omitempty"`

	// The reason code for the subscription date extension
	//
	// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
	ExtendReasonCode ExtendReasonCode `json:"extendReasonCode,omitempty"`

	// A string that contains a unique identifier you provide to track each subscription-renewal-date extension request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/requestidentifier
	RequestIdentifier string `json:"requestIdentifier,omitempty"`
}

// ExtendRenewalDateResponse is a response that indicates whether an individual renewal-date extension succeeded, and related details.
//
// https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldateresponse
type ExtendRenewalDateResponse struct {
	// The original transaction identifier of a purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/originaltransactionid
	OriginalTransactionId string `json:"originalTransactionId,omitempty"`

	// The unique identifier of subscription-purchase events across devices, including renewals.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/weborderlineitemid
	WebOrderLineItemId string `json:"webOrderLineItemId,omitempty"`

	// A Boolean value that indicates whether the subscription-renewal-date extension succeeded.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/success
	Success bool `json:"success,omitempty"`

	// The new subscription expiration date for a subscription-renewal extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/effectivedate
	EffectiveDate int64 `json:"effectiveDate,omitempty"`
}

// MassExtendRenewalDateRequest is the request body that contains subscription-renewal-extension data to apply for all eligible active subscribers.
//
// https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldaterequest
type MassExtendRenewalDateRequest struct {
	// The number of days to extend the subscription renewal date.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/extendbydays
	ExtendByDays int32 `json:"extendByDays,omitempty"`

	// The reason code for the subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
	ExtendReasonCode ExtendReasonCode `json:"extendReasonCode,omitempty"`

	// A string that contains a unique identifier you provide to track each subscription-renewal-date extension request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/requestidentifier
	RequestIdentifier string `json:"requestIdentifier,omitempty"`

	// A list of storefront country codes you provide to limit the storefronts for a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/storefrontcountrycodes
	StorefrontCountryCodes []string `json:"storefrontCountryCodes,omitempty"`

	// The unique identifier for the product, that you create in App Store Connect.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/productid
	ProductId string `json:"productId,omitempty"`
}

// MassExtendRenewalDateResponse is a response that indicates the server successfully received the subscription-renewal-date extension request.
//
// https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldateresponse
type MassExtendRenewalDateResponse struct {
	// A string that contains a unique identifier you provide to track each subscription-renewal-date extension request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/requestidentifier
	RequestIdentifier string `json:"requestIdentifier,omitempty"`
}

// MassExtendRenewalDateStatusResponse is a response that indicates the current status of a request to extend the subscription renewal date to all eligible subscribers.
//
// https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldatestatusresponse
type MassExtendRenewalDateStatusResponse struct {
	// A string that contains a unique identifier you provide to track each subscription-renewal-date extension request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/requestidentifier
	RequestIdentifier string `json:"requestIdentifier,omitempty"`

	// A Boolean value that indicates whether the App Store completed the request to extend a subscription renewal date to active subscribers.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/complete
	Complete bool `json:"complete,omitempty"`

	// The UNIX time, in milliseconds, that the App Store completes a request to extend a subscription renewal date for eligible subscribers.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/completedate
	CompleteDate int64 `json:"completeDate,omitempty"`

	// The count of subscriptions that successfully receive a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/succeededcount
	SucceededCount int64 `json:"succeededCount,omitempty"`

	// The count of subscriptions that fail to receive a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/failedcount
	FailedCount int64 `json:"failedCount,omitempty"`
}

// ConsumptionRequest is the request body that contains consumption information for an In-App Purchase.
//
// https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequest struct {
	// A Boolean value that indicates whether the customer consented to provide consumption data to the App Store.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/customerconsented
	CustomerConsented bool `json:"customerConsented"`

	// A Boolean value that indicates whether you provided, prior to its purchase, a free sample or trial of the content, or information about its functionality.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/samplecontentprovided
	SampleContentProvided bool `json:"sampleContentProvided"`

	// A value that indicates whether the app successfully delivered an in-app purchase that works properly.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/deliverystatus
	DeliveryStatus DeliveryStatus `json:"deliveryStatus"`

	// See deliveryStatus
	RawDeliveryStatus int32 `json:"rawDeliveryStatus"`

	// An integer that indicates the percentage, in milliunits, of the In-App Purchase the customer consumed.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/consumptionpercentage
	ConsumptionPercentage int32 `json:"consumptionPercentage,omitempty"`

	// A value that indicates your preferred outcome for the refund request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/refundpreference
	RefundPreference RefundPreference `json:"refundPreference,omitempty"`

	// See refundPreference
	RawRefundPreference int32 `json:"rawRefundPreference,omitempty"`
}

// UpdateAppAccountTokenRequest is the request body that contains an app account token value.
//
// https://developer.apple.com/documentation/appstoreserverapi/updateappaccounttokenrequest
type UpdateAppAccountTokenRequest struct {
	// The UUID that an app optionally generates to map a customer's in-app purchase with its resulting App Store transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/appaccounttoken
	AppAccountToken string `json:"appAccountToken"`
}

// SendTestNotificationResponse is a response that contains the test notification token.
//
// https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	// A unique identifier for a notification test that the App Store server sends to your server.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/testnotificationtoken
	TestNotificationToken string `json:"testNotificationToken,omitempty"`
}

// CheckTestNotificationResponse is a response that contains the contents of the test notification sent by the App Store server and the result from your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/checktestnotificationresponse
type CheckTestNotificationResponse struct {
	// A cryptographically signed payload, in JSON Web Signature (JWS) format, containing the response body for a version 2 notification.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/signedpayload
	SignedPayload string `json:"signedPayload,omitempty"`

	// An array of information the App Store server records for its attempts to send the TEST notification to your server. The array may contain a maximum of six sendAttemptItem objects.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
	SendAttempts []SendAttemptItem `json:"sendAttempts,omitempty"`
}

// AppTransactionInfoResponse is a response that contains signed app transaction information for a customer.
//
// https://developer.apple.com/documentation/appstoreserverapi/apptransactioninforesponse
type AppTransactionInfoResponse struct {
	// A customer’s app transaction information, signed by Apple, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwsapptransaction
	SignedAppTransactionInfo string `json:"signedAppTransactionInfo,omitempty"`
}

// JWSTransactionDecodedPayload is a decoded payload containing transaction information.
//
// https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type JWSTransactionDecodedPayload struct {
	// The original transaction identifier of a purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/originaltransactionid
	OriginalTransactionId string `json:"originalTransactionId,omitempty"`

	// The unique identifier for a transaction such as an in-app purchase, restored in-app purchase, or subscription renewal.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/transactionid
	TransactionId string `json:"transactionId,omitempty"`

	// The unique identifier of subscription-purchase events across devices, including renewals.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/weborderlineitemid
	WebOrderLineItemId string `json:"webOrderLineItemId,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The unique identifier for the product, that you create in App Store Connect.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/productid
	ProductId string `json:"productId,omitempty"`

	// The identifier of the subscription group that the subscription belongs to.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifier
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier,omitempty"`

	// The time that the App Store charged the user's account for an in-app purchase, a restored in-app purchase, a subscription, or a subscription renewal after a lapse.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/purchasedate
	PurchaseDate int64 `json:"purchaseDate,omitempty"`

	// The purchase date of the transaction associated with the original transaction identifier.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/originalpurchasedate
	OriginalPurchaseDate int64 `json:"originalPurchaseDate,omitempty"`

	// The UNIX time, in milliseconds, an auto-renewable subscription expires or renews.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/expiresdate
	ExpiresDate int64 `json:"expiresDate,omitempty"`

	// The number of consumable products purchased.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/quantity
	Quantity int32 `json:"quantity,omitempty"`

	// The type of the in-app purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/type
	Type Type `json:"type,omitempty"`

	// See type
	RawType string `json:"rawType,omitempty"`

	// The UUID that an app optionally generates to map a customer's in-app purchase with its resulting App Store transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/appaccounttoken
	AppAccountToken string `json:"appAccountToken,omitempty"`

	// A string that describes whether the transaction was purchased by the user, or is available to them through Family Sharing.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/inappownershiptype
	InAppOwnershipType InAppOwnershipType `json:"inAppOwnershipType,omitempty"`

	// See inAppOwnershipType
	RawInAppOwnershipType string `json:"rawInAppOwnershipType,omitempty"`

	// The UNIX time, in milliseconds, that the App Store signed the JSON Web Signature data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/signeddate
	SignedDate int64 `json:"signedDate,omitempty"`

	// The reason that the App Store refunded the transaction or revoked it from Family Sharing.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/revocationreason
	RevocationReason RevocationReason `json:"revocationReason,omitempty"`

	// See revocationReason
	RawRevocationReason int32 `json:"rawRevocationReason,omitempty"`

	// The UNIX time, in milliseconds, that Apple Support refunded a transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/revocationdate
	RevocationDate int64 `json:"revocationDate,omitempty"`

	// The Boolean value that indicates whether the user upgraded to another subscription.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/isupgraded
	IsUpgraded bool `json:"isUpgraded,omitempty"`

	// A value that represents the promotional offer type.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offertype
	OfferType OfferType `json:"offerType,omitempty"`

	// See offerType
	RawOfferType int32 `json:"rawOfferType,omitempty"`

	// The identifier that contains the offer code or the promotional offer identifier.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offeridentifier
	OfferIdentifier string `json:"offerIdentifier,omitempty"`

	// The server environment, either sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// The three-letter code that represents the country or region associated with the App Store storefront for the purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/storefront
	Storefront string `json:"storefront,omitempty"`

	// An Apple-defined value that uniquely identifies the App Store storefront associated with the purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/storefrontid
	StorefrontId string `json:"storefrontId,omitempty"`

	// The reason for the purchase transaction, which indicates whether it's a customer's purchase or a renewal for an auto-renewable subscription that the system initiates.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/transactionreason
	TransactionReason TransactionReason `json:"transactionReason,omitempty"`

	// See transactionReason
	RawTransactionReason string `json:"rawTransactionReason,omitempty"`

	// The three-letter ISO 4217 currency code for the price of the product.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/currency
	Currency string `json:"currency,omitempty"`

	// The price, in milliunits, of the in-app purchase or subscription offer that you configured in App Store Connect.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/price
	Price int64 `json:"price,omitempty"`

	// The payment mode you configure for the offer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offerdiscounttype
	OfferDiscountType OfferDiscountType `json:"offerDiscountType,omitempty"`

	// See offerDiscountType
	RawOfferDiscountType string `json:"rawOfferDiscountType,omitempty"`

	// The unique identifier of the app download transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/appTransactionId
	AppTransactionId string `json:"appTransactionId,omitempty"`

	// The duration of the offer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offerPeriod
	OfferPeriod string `json:"offerPeriod,omitempty"`

	// The type of the refund or revocation that applies to the transaction.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/revocationtype
	RevocationType RevocationType `json:"revocationType,omitempty"`

	// See revocationType
	RawRevocationType string `json:"rawRevocationType,omitempty"`

	// The percentage, in milliunits, of the transaction that the App Store has refunded or revoked.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/revocationpercentage
	RevocationPercentage int32 `json:"revocationPercentage,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *JWSTransactionDecodedPayload) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with all fields as any or specific types
	type Alias JWSTransactionDecodedPayload
	aux := &struct {
		Type               any `json:"type"`
		InAppOwnershipType any `json:"inAppOwnershipType"`
		RevocationReason   any `json:"revocationReason"`
		OfferType          any `json:"offerType"`
		Environment        any `json:"environment"`
		TransactionReason  any `json:"transactionReason"`
		OfferDiscountType  any `json:"offerDiscountType"`
		RevocationType     any `json:"revocationType"`
		// Floating point timestamps
		PurchaseDate         any `json:"purchaseDate"`
		OriginalPurchaseDate any `json:"originalPurchaseDate"`
		ExpiresDate          any `json:"expiresDate"`
		SignedDate           any `json:"signedDate"`
		RevocationDate       any `json:"revocationDate"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Type, &t.Type, &t.RawType, t.Type.Values())
	internal.UnmarshalStringEnum(aux.InAppOwnershipType, &t.InAppOwnershipType, &t.RawInAppOwnershipType, t.InAppOwnershipType.Values())
	internal.UnmarshalIntEnum(aux.RevocationReason, &t.RevocationReason, &t.RawRevocationReason)
	internal.UnmarshalIntEnum(aux.OfferType, &t.OfferType, &t.RawOfferType)
	internal.UnmarshalStringEnum(aux.Environment, &t.Environment, &t.RawEnvironment, t.Environment.Values())
	internal.UnmarshalStringEnum(aux.TransactionReason, &t.TransactionReason, &t.RawTransactionReason, t.TransactionReason.Values())
	internal.UnmarshalStringEnum(aux.OfferDiscountType, &t.OfferDiscountType, &t.RawOfferDiscountType, t.OfferDiscountType.Values())
	internal.UnmarshalStringEnum(aux.RevocationType, &t.RevocationType, &t.RawRevocationType, t.RevocationType.Values())
	internal.UnmarshalTimestamp(aux.PurchaseDate, &t.PurchaseDate)
	internal.UnmarshalTimestamp(aux.OriginalPurchaseDate, &t.OriginalPurchaseDate)
	internal.UnmarshalTimestamp(aux.ExpiresDate, &t.ExpiresDate)
	internal.UnmarshalTimestamp(aux.SignedDate, &t.SignedDate)
	internal.UnmarshalTimestamp(aux.RevocationDate, &t.RevocationDate)

	return nil
}

// JWSRenewalInfoDecodedPayload is a decoded payload containing subscription renewal information for an auto-renewable subscription.
//
// https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	// The reason the subscription expired.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/expirationintent
	ExpirationIntent ExpirationIntent `json:"expirationIntent,omitempty"`

	// See expirationIntent
	RawExpirationIntent int32 `json:"rawExpirationIntent,omitempty"`

	// The original transaction identifier of a purchase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/originaltransactionid
	OriginalTransactionId string `json:"originalTransactionId,omitempty"`

	// The product identifier of the product that will renew at the next billing period.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/autorenewproductid
	AutoRenewProductId string `json:"autoRenewProductId,omitempty"`

	// The unique identifier for the product, that you create in App Store Connect.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/productid
	ProductId string `json:"productId,omitempty"`

	// The renewal status of the auto-renewable subscription.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/autorenewstatus
	AutoRenewStatus AutoRenewStatus `json:"autoRenewStatus,omitempty"`

	// See autoRenewStatus
	RawAutoRenewStatus int32 `json:"rawAutoRenewStatus,omitempty"`

	// A Boolean value that indicates whether the App Store is attempting to automatically renew an expired subscription.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/isinbillingretryperiod
	IsInBillingRetryPeriod bool `json:"isInBillingRetryPeriod,omitempty"`

	// The status that indicates whether the auto-renewable subscription is subject to a price increase.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/priceincreasestatus
	PriceIncreaseStatus PriceIncreaseStatus `json:"priceIncreaseStatus,omitempty"`

	// See priceIncreaseStatus
	RawPriceIncreaseStatus int32 `json:"rawPriceIncreaseStatus,omitempty"`

	// The time when the billing grace period for subscription renewals expires.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/graceperiodexpiresdate
	GracePeriodExpiresDate int64 `json:"gracePeriodExpiresDate,omitempty"`

	// The type of subscription offer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offertype
	OfferType OfferType `json:"offerType,omitempty"`

	// See offerType
	RawOfferType int32 `json:"rawOfferType,omitempty"`

	// The offer code or the promotional offer identifier.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offeridentifier
	OfferIdentifier string `json:"offerIdentifier,omitempty"`

	// The UNIX time, in milliseconds, that the App Store signed the JSON Web Signature data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/signeddate
	SignedDate int64 `json:"signedDate,omitempty"`

	// The server environment, either sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// The earliest start date of a subscription in a series of auto-renewable subscription purchases that ignores all lapses of paid service shorter than 60 days.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/recentsubscriptionstartdate
	RecentSubscriptionStartDate int64 `json:"recentSubscriptionStartDate,omitempty"`

	// The UNIX time, in milliseconds, that the most recent auto-renewable subscription purchase expires.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/renewaldate
	RenewalDate int64 `json:"renewalDate,omitempty"`

	// The currency code for the renewalPrice of the subscription.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/currency
	Currency string `json:"currency,omitempty"`

	// The renewal price, in milliunits, of the auto-renewable subscription that renews at the next billing period.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/renewalprice
	RenewalPrice int64 `json:"renewalPrice,omitempty"`

	// The payment mode you configure for the offer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offerdiscounttype
	OfferDiscountType OfferDiscountType `json:"offerDiscountType,omitempty"`

	// See offerDiscountType
	RawOfferDiscountType string `json:"rawOfferDiscountType,omitempty"`

	// An array of win-back offer identifiers that a customer is eligible to redeem, which sorts the identifiers to present the better offers first.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/eligiblewinbackofferids
	EligibleWinBackOfferIds []string `json:"eligibleWinBackOfferIds,omitempty"`

	// The UUID that an app optionally generates to map a customer's in-app purchase with its resulting App Store transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/appaccounttoken
	AppAccountToken string `json:"appAccountToken,omitempty"`

	// The unique identifier of the app download transaction.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/appTransactionId
	AppTransactionId string `json:"appTransactionId,omitempty"`

	// The duration of the offer.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/offerPeriod
	OfferPeriod string `json:"offerPeriod,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *JWSRenewalInfoDecodedPayload) UnmarshalJSON(data []byte) error {
	type Alias JWSRenewalInfoDecodedPayload
	aux := &struct {
		ExpirationIntent    any `json:"expirationIntent"`
		AutoRenewStatus     any `json:"autoRenewStatus"`
		PriceIncreaseStatus any `json:"priceIncreaseStatus"`
		OfferType           any `json:"offerType"`
		Environment         any `json:"environment"`
		OfferDiscountType   any `json:"offerDiscountType"`
		// Floating point timestamps
		GracePeriodExpiresDate      any `json:"gracePeriodExpiresDate"`
		SignedDate                  any `json:"signedDate"`
		RecentSubscriptionStartDate any `json:"recentSubscriptionStartDate"`
		RenewalDate                 any `json:"renewalDate"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalIntEnum(aux.ExpirationIntent, &t.ExpirationIntent, &t.RawExpirationIntent)
	internal.UnmarshalIntEnum(aux.AutoRenewStatus, &t.AutoRenewStatus, &t.RawAutoRenewStatus)
	internal.UnmarshalIntEnum(aux.PriceIncreaseStatus, &t.PriceIncreaseStatus, &t.RawPriceIncreaseStatus)
	internal.UnmarshalIntEnum(aux.OfferType, &t.OfferType, &t.RawOfferType)
	internal.UnmarshalStringEnum(aux.Environment, &t.Environment, &t.RawEnvironment, t.Environment.Values())
	internal.UnmarshalStringEnum(aux.OfferDiscountType, &t.OfferDiscountType, &t.RawOfferDiscountType, t.OfferDiscountType.Values())
	internal.UnmarshalTimestamp(aux.GracePeriodExpiresDate, &t.GracePeriodExpiresDate)
	internal.UnmarshalTimestamp(aux.SignedDate, &t.SignedDate)
	internal.UnmarshalTimestamp(aux.RecentSubscriptionStartDate, &t.RecentSubscriptionStartDate)
	internal.UnmarshalTimestamp(aux.RenewalDate, &t.RenewalDate)

	return nil
}

// AppTransaction is information that represents the customer’s purchase of the app, cryptographically signed by the App Store.
//
// https://developer.apple.com/documentation/storekit/apptransaction
type AppTransaction struct {
	// The server environment that signs the app transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/environment
	ReceiptType Environment `json:"receiptType,omitempty"`

	// See receiptType
	RawReceiptType string `json:"rawReceiptType,omitempty"`

	// The unique identifier the App Store uses to identify the app.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/appid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier that the app transaction applies to.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The app version that the app transaction applies to.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/appversion
	ApplicationVersion string `json:"applicationVersion,omitempty"`

	// The version external identifier of the app
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/appversionid
	VersionExternalIdentifier int64 `json:"versionExternalIdentifier,omitempty"`

	// The date that the App Store signed the JWS app transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/signeddate
	ReceiptCreationDate int64 `json:"receiptCreationDate,omitempty"`

	// The date the user originally purchased the app from the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/originalpurchasedate
	OriginalPurchaseDate int64 `json:"originalPurchaseDate,omitempty"`

	// The app version that the user originally purchased from the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/originalappversion
	OriginalApplicationVersion string `json:"originalApplicationVersion,omitempty"`

	// The Base64 device verification value to use to verify whether the app transaction belongs to the device.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/deviceverification
	DeviceVerification string `json:"deviceVerification,omitempty"`

	// The UUID used to compute the device verification value.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/deviceverificationnonce
	DeviceVerificationNonce string `json:"deviceVerificationNonce,omitempty"`

	// The date the customer placed an order for the app before it's available in the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/preorderdate
	PreorderDate int64 `json:"preorderDate,omitempty"`

	// The unique identifier of the app download transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/apptransactionid
	AppTransactionId string `json:"appTransactionId,omitempty"`

	// The platform on which the customer originally purchased the app.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/originalplatform
	OriginalPlatform PurchasePlatform `json:"originalPlatform,omitempty"`

	// See originalPlatform
	RawOriginalPlatform string `json:"rawOriginalPlatform,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *AppTransaction) UnmarshalJSON(data []byte) error {
	type Alias AppTransaction
	aux := &struct {
		ReceiptType      any `json:"receiptType"`
		OriginalPlatform any `json:"originalPlatform"`
		// Floating point timestamps
		ReceiptCreationDate  any `json:"receiptCreationDate"`
		OriginalPurchaseDate any `json:"originalPurchaseDate"`
		PreorderDate         any `json:"preorderDate"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.ReceiptType, &t.ReceiptType, &t.RawReceiptType, t.ReceiptType.Values())
	internal.UnmarshalStringEnum(aux.OriginalPlatform, &t.OriginalPlatform, &t.RawOriginalPlatform, t.OriginalPlatform.Values())
	internal.UnmarshalTimestamp(aux.ReceiptCreationDate, &t.ReceiptCreationDate)
	internal.UnmarshalTimestamp(aux.OriginalPurchaseDate, &t.OriginalPurchaseDate)
	internal.UnmarshalTimestamp(aux.PreorderDate, &t.PreorderDate)

	return nil
}
