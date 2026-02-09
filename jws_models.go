package appstore

import (
	"encoding/json"

	"github.com/laishere/app-store-server-library-go/internal"
)

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

// UnmarshalJSON custom unmarshaler for JWSTransactionDecodedPayload
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

	// Handle Type
	internal.UnmarshalStringEnum(aux.Type, &t.Type, &t.RawType)

	// Handle InAppOwnershipType
	internal.UnmarshalStringEnum(aux.InAppOwnershipType, &t.InAppOwnershipType, &t.RawInAppOwnershipType)

	// Handle RevocationReason
	internal.UnmarshalIntEnum(aux.RevocationReason, &t.RevocationReason, &t.RawRevocationReason)

	// Handle OfferType
	internal.UnmarshalIntEnum(aux.OfferType, &t.OfferType, &t.RawOfferType)

	// Handle Environment
	internal.UnmarshalStringEnum(aux.Environment, &t.Environment, &t.RawEnvironment)

	// Handle TransactionReason
	internal.UnmarshalStringEnum(aux.TransactionReason, &t.TransactionReason, &t.RawTransactionReason)

	// Handle OfferDiscountType
	internal.UnmarshalStringEnum(aux.OfferDiscountType, &t.OfferDiscountType, &t.RawOfferDiscountType)

	// Handle RevocationType
	internal.UnmarshalStringEnum(aux.RevocationType, &t.RevocationType, &t.RawRevocationType)

	// Handle floating point timestamps
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

// UnmarshalJSON custom unmarshaler for JWSRenewalInfoDecodedPayload
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

	// Handle ExpirationIntent
	internal.UnmarshalIntEnum(aux.ExpirationIntent, &t.ExpirationIntent, &t.RawExpirationIntent)

	// Handle AutoRenewStatus
	internal.UnmarshalIntEnum(aux.AutoRenewStatus, &t.AutoRenewStatus, &t.RawAutoRenewStatus)

	// Handle PriceIncreaseStatus
	internal.UnmarshalIntEnum(aux.PriceIncreaseStatus, &t.PriceIncreaseStatus, &t.RawPriceIncreaseStatus)

	// Handle OfferType
	internal.UnmarshalIntEnum(aux.OfferType, &t.OfferType, &t.RawOfferType)

	// Handle Environment
	internal.UnmarshalStringEnum(aux.Environment, &t.Environment, &t.RawEnvironment)

	// Handle OfferDiscountType
	internal.UnmarshalStringEnum(aux.OfferDiscountType, &t.OfferDiscountType, &t.RawOfferDiscountType)

	// Handle floating point timestamps
	internal.UnmarshalTimestamp(aux.GracePeriodExpiresDate, &t.GracePeriodExpiresDate)
	internal.UnmarshalTimestamp(aux.SignedDate, &t.SignedDate)
	internal.UnmarshalTimestamp(aux.RecentSubscriptionStartDate, &t.RecentSubscriptionStartDate)
	internal.UnmarshalTimestamp(aux.RenewalDate, &t.RenewalDate)

	return nil
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

// UnmarshalJSON custom unmarshaler for Data
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

	// Handle Environment
	internal.UnmarshalStringEnum(aux.Environment, &d.Environment, &d.RawEnvironment)

	// Handle Status
	internal.UnmarshalIntEnum(aux.Status, &d.Status, &d.RawStatus)

	// Handle ConsumptionRequestReason
	internal.UnmarshalStringEnum(aux.ConsumptionRequestReason, &d.ConsumptionRequestReason, &d.RawConsumptionRequestReason)

	return nil
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

// UnmarshalJSON custom unmarshaler for AppData to populate RawEnvironment
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

	// Handle Environment - populate both parsed enum and raw string
	internal.UnmarshalStringEnum(aux.Environment, &a.Environment, &a.RawEnvironment)

	return nil
}

// AppTransaction is information that represents the customer’s purchase of the app, cryptographically signed by the App Store.
//
// https://developer.apple.com/documentation/storekit/apptransaction
type AppTransaction struct {
	// The server environment that signs the app transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3963901-environment
	ReceiptType Environment `json:"receiptType,omitempty"`

	// See receiptType
	RawReceiptType string `json:"rawReceiptType,omitempty"`

	// The unique identifier the App Store uses to identify the app.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954436-appid
	AppAppleId int64 `json:"appAppleId,omitempty"`

	// The bundle identifier that the app transaction applies to.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954439-bundleid
	BundleId string `json:"bundleId,omitempty"`

	// The app version that the app transaction applies to.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954437-appversion
	ApplicationVersion string `json:"applicationVersion,omitempty"`

	// The version external identifier of the app
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954438-appversionid
	VersionExternalIdentifier int64 `json:"versionExternalIdentifier,omitempty"`

	// The date that the App Store signed the JWS app transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954449-signeddate
	ReceiptCreationDate int64 `json:"receiptCreationDate,omitempty"`

	// The date the user originally purchased the app from the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954448-originalpurchasedate
	OriginalPurchaseDate int64 `json:"originalPurchaseDate,omitempty"`

	// The app version that the user originally purchased from the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954447-originalappversion
	OriginalApplicationVersion string `json:"originalApplicationVersion,omitempty"`

	// The Base64 device verification value to use to verify whether the app transaction belongs to the device.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954441-deviceverification
	DeviceVerification string `json:"deviceVerification,omitempty"`

	// The UUID used to compute the device verification value.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/3954442-deviceverificationnonce
	DeviceVerificationNonce string `json:"deviceVerificationNonce,omitempty"`

	// The date the customer placed an order for the app before it's available in the App Store.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/4013175-preorderdate
	PreorderDate int64 `json:"preorderDate,omitempty"`

	// The unique identifier of the app download transaction.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/apptransactionid
	AppTransactionId string `json:"appTransactionId,omitempty"`

	// The platform on which the customer originally purchased the app.
	//
	// https://developer.apple.com/documentation/storekit/apptransaction/originalplatform-4mogz
	OriginalPlatform PurchasePlatform `json:"originalPlatform,omitempty"`

	// See originalPlatform
	RawOriginalPlatform string `json:"rawOriginalPlatform,omitempty"`
}

// UnmarshalJSON custom unmarshaler for AppTransaction
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

	// Handle ReceiptType (Environment)
	internal.UnmarshalStringEnum(aux.ReceiptType, &t.ReceiptType, &t.RawReceiptType)

	// Handle OriginalPlatform
	internal.UnmarshalStringEnum(aux.OriginalPlatform, &t.OriginalPlatform, &t.RawOriginalPlatform)

	// Handle floating point timestamps
	internal.UnmarshalTimestamp(aux.ReceiptCreationDate, &t.ReceiptCreationDate)
	internal.UnmarshalTimestamp(aux.OriginalPurchaseDate, &t.OriginalPurchaseDate)
	internal.UnmarshalTimestamp(aux.PreorderDate, &t.PreorderDate)

	return nil
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
