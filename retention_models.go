package appstore

import (
	"encoding/json"

	"github.com/laishere/app-store-server-library-go/internal"
)

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

// UnmarshalJSON implements json.Unmarshaler.
func (g *GetImageListResponseItem) UnmarshalJSON(data []byte) error {
	type Alias GetImageListResponseItem
	aux := &struct {
		ImageState any `json:"imageState"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.ImageState, &g.ImageState, &g.RawImageState, g.ImageState.Values())

	return nil
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

// UnmarshalJSON implements json.Unmarshaler.
func (g *GetMessageListResponseItem) UnmarshalJSON(data []byte) error {
	type Alias GetMessageListResponseItem
	aux := &struct {
		MessageState any `json:"messageState"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.MessageState, &g.MessageState, &g.RawMessageState, g.MessageState.Values())

	return nil
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
	Timestamp *Timestamp `json:"timestamp,omitempty"`

	// A string that identifies the private key you use to generate the signature.
	KeyId *string `json:"keyId,omitempty"`

	// The subscription offer identifier that you set up in App Store Connect.
	OfferIdentifier *string `json:"offerIdentifier,omitempty"`

	// A UUID that you provide to associate with the transaction if the customer accepts the promotional offer.
	AppAccountToken *string `json:"appAccountToken,omitempty"`
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
	SignedDate Timestamp `json:"signedDate"`

	// The server environment, either sandbox or production.
	//
	// https://developer.apple.com/documentation/retentionmessaging/environment
	Environment Environment `json:"environment"`

	// See environment
	RawEnvironment string `json:"rawEnvironment"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DecodedRealtimeRequestBody) UnmarshalJSON(data []byte) error {
	type Alias DecodedRealtimeRequestBody
	aux := &struct {
		Environment any `json:"environment"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	internal.UnmarshalStringEnum(aux.Environment, &d.Environment, &d.RawEnvironment, d.Environment.Values())

	return nil
}
