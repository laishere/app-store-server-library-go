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
	Environment *Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// An array of in-app purchase transactions for the customer, signed by Apple, in JSON Web Signature format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactions []string `json:"signedTransactions,omitempty"`
}

// UnmarshalJSON custom unmarshaler for HistoryResponse to populate RawEnvironment
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

	// Handle Environment - populate both parsed enum and raw string
	var env Environment
	internal.UnmarshalStringEnum(aux.Environment, &env, &h.RawEnvironment)
	if h.RawEnvironment != "" {
		switch env {
		case ENVIRONMENT_PRODUCTION, ENVIRONMENT_SANDBOX, ENVIRONMENT_XCODE, ENVIRONMENT_LOCAL_TESTING:
			h.Environment = &env
		default:
			h.Environment = nil
		}
	}

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

// UnmarshalJSON custom unmarshaler for StatusResponse to populate RawEnvironment
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

	// Handle Environment - populate both parsed enum and raw string
	internal.UnmarshalStringEnum(aux.Environment, &s.Environment, &s.RawEnvironment)

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

// DefaultConfigurationRequest is the request body that contains the default configuration information.
//
// https://developer.apple.com/documentation/retentionmessaging/defaultconfigurationrequest
type DefaultConfigurationRequest struct {
	// The message identifier of the message to configure as a default message.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageidentifier
	MessageIdentifier string `json:"messageIdentifier,omitempty"`
}

// GetImageListResponse is a response that contains status information for all images.
//
// https://developer.apple.com/documentation/retentionmessaging/getimagelistresponse
type GetImageListResponse struct {
	// An array of all image identifiers and their image state.
	//
	// https://developer.apple.com/documentation/retentionmessaging/getimagelistresponseitem
	ImageIdentifiers []GetImageListResponseItem `json:"imageIdentifiers,omitempty"`
}

// GetImageListResponseItem is an image identifier and state information for an image.
//
// https://developer.apple.com/documentation/retentionmessaging/getimagelistresponseitem
type GetImageListResponseItem struct {
	// The identifier of the image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/imageidentifier
	ImageIdentifier string `json:"imageIdentifier,omitempty"`

	// The current state of the image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/imagestate
	ImageState ImageState `json:"imageState,omitempty"`

	// See imageState
	RawImageState string `json:"rawImageState,omitempty"`
}

// GetMessageListResponse is a response that contains status information for all messages.
//
// https://developer.apple.com/documentation/retentionmessaging/getmessagelistresponse
type GetMessageListResponse struct {
	// An array of all message identifiers and their message state.
	//
	// https://developer.apple.com/documentation/retentionmessaging/getmessagelistresponseitem
	MessageIdentifiers []GetMessageListResponseItem `json:"messageIdentifiers,omitempty"`
}

// GetMessageListResponseItem is a message identifier and status information for a message.
//
// https://developer.apple.com/documentation/retentionmessaging/getmessagelistresponseitem
type GetMessageListResponseItem struct {
	// The identifier of the message.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageidentifier
	MessageIdentifier string `json:"messageIdentifier,omitempty"`

	// The current state of the message.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageState
	MessageState MessageState `json:"messageState,omitempty"`

	// See messageState
	RawMessageState string `json:"rawMessageState,omitempty"`
}

// UploadMessageRequestBody is the request body for uploading a message, which includes the message text and an optional image reference.
//
// https://developer.apple.com/documentation/retentionmessaging/uploadmessagerequestbody
type UploadMessageRequestBody struct {
	// The header text of the retention message that the system displays to customers.
	//
	// https://developer.apple.com/documentation/retentionmessaging/header
	Header string `json:"header"`

	// The body text of the retention message that the system displays to customers.
	//
	// https://developer.apple.com/documentation/retentionmessaging/body
	Body string `json:"body"`

	// The optional image identifier and its alternative text to appear as part of a text-based message with an image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/uploadmessageimage
	Image *UploadMessageImage `json:"image,omitempty"`
}

// UploadMessageImage is the definition of an image with its alternative text.
//
// https://developer.apple.com/documentation/retentionmessaging/uploadmessageimage
type UploadMessageImage struct {
	// The unique identifier of an image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/imageidentifier
	ImageIdentifier string `json:"imageIdentifier"`

	// The alternative text you provide for the corresponding image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/alttext
	AltText string `json:"altText"`
}

// RealtimeResponseBody is the response you provide to choose, in real time, a retention message the system displays to the customer.
//
// https://developer.apple.com/documentation/retentionmessaging/realtimeresponsebody
type RealtimeResponseBody struct {
	// A retention message that's text-based and can include an optional image.
	//
	// https://developer.apple.com/documentation/retentionmessaging/message
	Message *Message `json:"message,omitempty"`

	// A retention message with a switch-plan option.
	//
	// https://developer.apple.com/documentation/retentionmessaging/alternateproduct
	AlternateProduct *AlternateProduct `json:"alternateProduct,omitempty"`

	// A retention message that includes a promotional offer.
	//
	// https://developer.apple.com/documentation/retentionmessaging/promotionaloffer
	PromotionalOffer *PromotionalOffer `json:"promotionalOffer,omitempty"`
}

// Message is a message identifier you provide in a real-time response to your Get Retention Message endpoint.
//
// https://developer.apple.com/documentation/retentionmessaging/message
type Message struct {
	// The identifier of the message to display to the customer.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageidentifier
	MessageIdentifier *string `json:"messageIdentifier,omitempty"`
}

// AlternateProduct is a switch-plan message and product ID you provide in a real-time response to your Get Retention Message endpoint.
//
// https://developer.apple.com/documentation/retentionmessaging/alternateproduct
type AlternateProduct struct {
	// The message identifier of the text to display in the switch-plan retention message.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageidentifier
	MessageIdentifier *string `json:"messageIdentifier,omitempty"`

	// The product identifier of the subscription the retention message suggests for your customer to switch to.
	//
	// https://developer.apple.com/documentation/retentionmessaging/productid
	ProductId *string `json:"productId,omitempty"`
}

// PromotionalOffer is a promotional offer and message you provide in a real-time response to your Get Retention Message endpoint.
//
// https://developer.apple.com/documentation/retentionmessaging/promotionaloffer
type PromotionalOffer struct {
	// The identifier of the message to display to the customer, along with the promotional offer.
	//
	// https://developer.apple.com/documentation/retentionmessaging/messageidentifier
	MessageIdentifier *string `json:"messageIdentifier,omitempty"`

	// The promotional offer signature in V2 format.
	//
	// https://developer.apple.com/documentation/retentionmessaging/promotionaloffersignaturev2
	PromotionalOfferSignatureV2 *string `json:"promotionalOfferSignatureV2,omitempty"`

	// The promotional offer signature in V1 format.
	//
	// https://developer.apple.com/documentation/retentionmessaging/promotionaloffersignaturev1
	PromotionalOfferSignatureV1 *PromotionalOfferSignatureV1 `json:"promotionalOfferSignatureV1,omitempty"`
}

// PromotionalOfferSignatureV1 is the promotional offer signature you generate using an earlier signature version.
//
// https://developer.apple.com/documentation/retentionmessaging/promotionaloffersignaturev1
type PromotionalOfferSignatureV1 struct {
	// The Base64-encoded cryptographic signature you generate using the offer parameters.
	EncodedSignature *string `json:"encodedSignature,omitempty"`

	// The subscription's product identifier.
	//
	// https://developer.apple.com/documentation/retentionmessaging/productid
	ProductId *string `json:"productId,omitempty"`

	// A one-time-use UUID antireplay value you generate.
	Nonce *string `json:"nonce,omitempty"`

	// The UNIX time, in milliseconds, when you generate the signature.
	Timestamp *int64 `json:"timestamp,omitempty"`

	// A string that identifies the private key you use to generate the signature.
	KeyId *string `json:"keyId,omitempty"`

	// The subscription offer identifier that you set up in App Store Connect.
	OfferIdentifier *string `json:"offerIdentifier,omitempty"`

	// A UUID that you provide to associate with the transaction if the customer accepts the promotional offer.
	AppAccountToken *string `json:"appAccountToken,omitempty"`
}
