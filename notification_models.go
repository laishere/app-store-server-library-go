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

// UnmarshalJSON custom unmarshaler for ResponseBodyV2DecodedPayload
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

	// Handle NotificationType
	internal.UnmarshalStringEnum(aux.NotificationType, &t.NotificationType, &t.RawNotificationType)

	// Handle Subtype
	internal.UnmarshalStringEnum(aux.Subtype, &t.Subtype, &t.RawSubtype)

	// Handle floating point timestamps
	internal.UnmarshalTimestamp(aux.SignedDate, &t.SignedDate)

	return nil
}

// DecodedRealtimeRequestBody is the decoded request body the App Store sends to your server to request a real-time retention message.
//
// https://developer.apple.com/documentation/retentionmessaging/decodedrealtimerequestbody
type DecodedRealtimeRequestBody struct {
	// The original transaction identifier of the customer's subscription.
	//
	// https://developer.apple.com/documentation/retentionmessaging/originaltransactionid
	OriginalTransactionId string `json:"originalTransactionId"`

	// The unique identifier of the app in the App Store.
	//
	// https://developer.apple.com/documentation/retentionmessaging/appappleid
	AppAppleId int64 `json:"appAppleId"`

	// The unique identifier of the app in the App Store.
	//
	// https://developer.apple.com/documentation/retentionmessaging/productid
	ProductId string `json:"productId"`

	// The device's locale.
	//
	// https://developer.apple.com/documentation/retentionmessaging/locale
	UserLocale string `json:"userLocale"`

	// A UUID the App Store server creates to uniquely identify each request.
	//
	// https://developer.apple.com/documentation/retentionmessaging/requestidentifier
	RequestIdentifier string `json:"requestIdentifier"`

	// The UNIX time, in milliseconds, that the App Store signed the JSON Web Signature (JWS) data.
	//
	// https://developer.apple.com/documentation/retentionmessaging/signeddate
	SignedDate int64 `json:"signedDate"`

	// The server environment, either sandbox or production.
	//
	// https://developer.apple.com/documentation/retentionmessaging/environment
	Environment Environment `json:"environment"`

	// See environment
	RawEnvironment string `json:"rawEnvironment"`
}

// UnmarshalJSON custom unmarshaler for DecodedRealtimeRequestBody to populate RawEnvironment
func (d *DecodedRealtimeRequestBody) UnmarshalJSON(data []byte) error {
	type Alias DecodedRealtimeRequestBody
	aux := &struct {
		Environment any `json:"environment"`
		SignedDate  any `json:"signedDate"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle Environment - populate both parsed enum and raw string
	internal.UnmarshalStringEnum(aux.Environment, &d.Environment, &d.RawEnvironment)

	// Handle floating point timestamps
	internal.UnmarshalTimestamp(aux.SignedDate, &d.SignedDate)

	return nil
}
