package appstore

import (
	"encoding/json"

	"github.com/laishere/app-store-server-library-go/internal"
)

// ResponseBodyV2DecodedPayload is a decoded payload containing the version 2 notification data.
//
// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
type ResponseBodyV2DecodedPayload struct {
	// The in-app purchase event for which the App Store sends this version 2 notification.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
	NotificationType NotificationTypeV2 `json:"notificationType,omitempty"`

	// See notificationType
	RawNotificationType string `json:"rawNotificationType,omitempty"`

	// Additional information that identifies the notification event. The subtype field is present only for specific version 2 notifications.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/subtype
	Subtype Subtype `json:"subtype,omitempty"`

	// See subtype
	RawSubtype string `json:"rawSubtype,omitempty"`

	// A unique identifier for the notification.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/notificationuuid
	NotificationUUID string `json:"notificationUUID,omitempty"`

	// The object that contains the app metadata and signed renewal and transaction information.
	// The data, summary, and externalPurchaseToken fields are mutually exclusive. The payload contains only one of these fields.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/data
	Data *Data `json:"data,omitempty"`

	// A string that indicates the notification's App Store Server Notifications version number.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/version
	Version string `json:"version,omitempty"`

	// The UNIX time, in milliseconds, that the App Store signed the JSON Web Signature data.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/signeddate
	SignedDate int64 `json:"signedDate,omitempty"`

	// The summary data that appears when the App Store server completes your request to extend a subscription renewal date for eligible subscribers.
	// The data, summary, and externalPurchaseToken fields are mutually exclusive. The payload contains only one of these fields.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/summary
	Summary *Summary `json:"summary,omitempty"`

	// This field appears when the notificationType is EXTERNAL_PURCHASE_TOKEN.
	// The data, summary, and externalPurchaseToken fields are mutually exclusive. The payload contains only one of these fields.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/externalpurchasetoken
	ExternalPurchaseToken *ExternalPurchaseToken `json:"externalPurchaseToken,omitempty"`

	// The object that contains the app metadata and signed app transaction information. This field appears when the notificationType is RESCIND_CONSENT.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appdata
	AppData *AppData `json:"appData,omitempty"`
}

// Data is the app metadata and the signed renewal and transaction information.
//
// https://developer.apple.com/documentation/appstoreservernotifications/data
type Data struct {
	// The server environment that the notification applies to, either sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// The unique identifier of an app in the App Store.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The version of the build that identifies an iteration of the bundle.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/bundleversion
	BundleVersion string `json:"bundleVersion,omitempty"`

	// Transaction information signed by the App Store, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
	SignedTransactionInfo string `json:"signedTransactionInfo,omitempty"`

	// Subscription renewal information, signed by the App Store, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfo
	SignedRenewalInfo string `json:"signedRenewalInfo,omitempty"`

	// The status of an auto-renewable subscription as of the signedDate in the responseBodyV2DecodedPayload.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/status
	Status Status `json:"status,omitempty"`

	// See status
	RawStatus int32 `json:"rawStatus,omitempty"`

	// The reason the customer requested the refund.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/consumptionrequestreason
	ConsumptionRequestReason ConsumptionRequestReason `json:"consumptionRequestReason,omitempty"`

	// See consumptionRequestReason
	RawConsumptionRequestReason string `json:"rawConsumptionRequestReason,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Data) UnmarshalJSON(data []byte) error {
	type Alias Data
	aux := &struct {
		Environment              any `json:"environment"`
		Status                   any `json:"status"`
		ConsumptionRequestReason any `json:"consumptionRequestReason"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Environment, &d.Environment, &d.RawEnvironment, d.Environment.Values())
	internal.UnmarshalIntEnum(aux.Status, &d.Status, &d.RawStatus)
	internal.UnmarshalStringEnum(aux.ConsumptionRequestReason, &d.ConsumptionRequestReason, &d.RawConsumptionRequestReason, d.ConsumptionRequestReason.Values())

	return nil
}

// Summary is the payload data for a subscription-renewal-date extension notification.
//
// https://developer.apple.com/documentation/appstoreservernotifications/summary
type Summary struct {
	// The server environment that the notification applies to, either sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// The unique identifier of an app in the App Store.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The unique identifier for the product, that you create in App Store Connect.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/productid
	ProductId string `json:"productId,omitempty"`

	// A string that contains a unique identifier you provide to track each subscription-renewal-date extension request.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/requestidentifier
	RequestIdentifier string `json:"requestIdentifier,omitempty"`

	// A list of storefront country codes you provide to limit the storefronts for a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/storefrontcountrycodes
	StorefrontCountryCodes []string `json:"storefrontCountryCodes,omitempty"`

	// The count of subscriptions that successfully receive a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/succeededcount
	SucceededCount int64 `json:"succeededCount,omitempty"`

	// The count of subscriptions that fail to receive a subscription-renewal-date extension.
	//
	// https://developer.apple.com/documentation/appstoreserverapi/failedcount
	FailedCount int64 `json:"failedCount,omitempty"`
}

// ExternalPurchaseToken is the payload data that contains an external purchase token.
//
// https://developer.apple.com/documentation/appstoreservernotifications/externalpurchasetoken
type ExternalPurchaseToken struct {
	// The field of an external purchase token that uniquely identifies the token.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/externalpurchaseid
	ExternalPurchaseId string `json:"externalPurchaseId,omitempty"`

	// The field of an external purchase token that contains the UNIX date, in milliseconds, when the system created the token.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/tokencreationdate
	TokenCreationDate int64 `json:"tokenCreationDate,omitempty"`

	// The unique identifier of an app in the App Store.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier of an app.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/bundleid
	BundleId string `json:"bundleId,omitempty"`
}

// AppData is the object that contains the app metadata and signed app transaction information.
//
// https://developer.apple.com/documentation/appstoreservernotifications/appdata
type AppData struct {
	// The unique identifier of the app that the notification applies to.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/appappleid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier of the app.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The server environment that the notification applies to, either sandbox or production.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/environment
	Environment Environment `json:"environment,omitempty"`

	// See environment
	RawEnvironment string `json:"rawEnvironment,omitempty"`

	// App transaction information signed by the App Store, in JSON Web Signature (JWS) format.
	//
	// https://developer.apple.com/documentation/appstoreservernotifications/jwsapptransaction
	SignedAppTransactionInfo string `json:"signedAppTransactionInfo,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AppData) UnmarshalJSON(data []byte) error {
	type Alias AppData
	aux := &struct {
		Environment any `json:"environment"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Environment, &a.Environment, &a.RawEnvironment, a.Environment.Values())

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *ResponseBodyV2DecodedPayload) UnmarshalJSON(data []byte) error {
	type Alias ResponseBodyV2DecodedPayload
	aux := &struct {
		NotificationType any `json:"notificationType"`
		Subtype          any `json:"subtype"`
		SignedDate       any `json:"signedDate"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.NotificationType, &t.NotificationType, &t.RawNotificationType, t.NotificationType.Values())
	internal.UnmarshalStringEnum(aux.Subtype, &t.Subtype, &t.RawSubtype, t.Subtype.Values())
	internal.UnmarshalTimestamp(aux.SignedDate, &t.SignedDate)

	return nil
}
