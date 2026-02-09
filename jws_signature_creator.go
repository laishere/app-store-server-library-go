package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWSSignatureCreator creates signed JWS tokens for various App Store features.
// It encapsulates the signing key and associated metadata needed for creating signatures.
type JWSSignatureCreator struct {
	audience   string
	signingKey *ecdsa.PrivateKey
	keyID      string
	issuerID   string
	bundleID   string
}

// NewJWSSignatureCreator creates a new JWS signature creator.
// The signingKey should be a PEM-encoded PKCS#8 ECDSA private key.
func NewJWSSignatureCreator(audience string, signingKey []byte, keyID, issuerID, bundleID string) (*JWSSignatureCreator, error) {
	block, _ := pem.Decode(signingKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an ECDSA private key")
	}

	return &JWSSignatureCreator{
		audience:   audience,
		signingKey: ecdsaKey,
		keyID:      keyID,
		issuerID:   issuerID,
		bundleID:   bundleID,
	}, nil
}

func (s *JWSSignatureCreator) createSignature(featureSpecificClaims map[string]any) (string, error) {
	claims := jwt.MapClaims{}
	maps.Copy(claims, featureSpecificClaims)

	claims["bid"] = s.bundleID
	claims["iss"] = s.issuerID
	claims["aud"] = s.audience
	claims["iat"] = time.Now().Unix()
	claims["nonce"] = uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = s.keyID

	return token.SignedString(s.signingKey)
}

// PromotionalOfferV2SignatureCreator creates signatures for promotional offers.
// It is used to generate signed offer tokens for App Store promotional offers.
type PromotionalOfferV2SignatureCreator struct {
	*JWSSignatureCreator
}

// NewPromotionalOfferV2SignatureCreator creates a new promotional offer signature creator.
// The signingKey should be a PEM-encoded PKCS#8 ECDSA private key.
func NewPromotionalOfferV2SignatureCreator(signingKey []byte, keyID, issuerID, bundleID string) (*PromotionalOfferV2SignatureCreator, error) {
	base, err := NewJWSSignatureCreator("promotional-offer", signingKey, keyID, issuerID, bundleID)
	if err != nil {
		return nil, err
	}
	return &PromotionalOfferV2SignatureCreator{base}, nil
}

// CreateSignature creates a signed token for a promotional offer.
// The transactionID parameter is optional and may be nil.
func (s *PromotionalOfferV2SignatureCreator) CreateSignature(productID, offerIdentifier string, transactionID *string) (string, error) {
	if productID == "" {
		return "", errors.New("product_id cannot be null")
	}
	if offerIdentifier == "" {
		return "", errors.New("offer_identifier cannot be null")
	}
	featureSpecificClaims := map[string]any{
		"productId":       productID,
		"offerIdentifier": offerIdentifier,
	}
	if transactionID != nil {
		featureSpecificClaims["transactionId"] = *transactionID
	}
	return s.createSignature(featureSpecificClaims)
}

// IntroductoryOfferEligibilitySignatureCreator creates signatures for checking introductory offer eligibility.
// It is used to generate signed tokens for verifying if a user is eligible for an introductory offer.
type IntroductoryOfferEligibilitySignatureCreator struct {
	*JWSSignatureCreator
}

// NewIntroductoryOfferEligibilitySignatureCreator creates a new introductory offer eligibility signature creator.
// The signingKey should be a PEM-encoded PKCS#8 ECDSA private key.
func NewIntroductoryOfferEligibilitySignatureCreator(signingKey []byte, keyID, issuerID, bundleID string) (*IntroductoryOfferEligibilitySignatureCreator, error) {
	base, err := NewJWSSignatureCreator("introductory-offer-eligibility", signingKey, keyID, issuerID, bundleID)
	if err != nil {
		return nil, err
	}
	return &IntroductoryOfferEligibilitySignatureCreator{base}, nil
}

// CreateSignature creates a signed token to check introductory offer eligibility.
func (s *IntroductoryOfferEligibilitySignatureCreator) CreateSignature(productID string, allowIntroductoryOffer bool, transactionID string) (string, error) {
	if productID == "" {
		return "", errors.New("product_id cannot be null")
	}
	if transactionID == "" {
		return "", errors.New("transaction_id cannot be null")
	}
	featureSpecificClaims := map[string]any{
		"productId":              productID,
		"allowIntroductoryOffer": allowIntroductoryOffer,
		"transactionId":          transactionID,
	}
	return s.createSignature(featureSpecificClaims)
}

// AdvancedCommerceAPIInAppSignatureCreator creates signatures for Advanced Commerce API requests.
// It is used to generate signed tokens for advanced commerce in-app purchase requests.
type AdvancedCommerceAPIInAppSignatureCreator struct {
	*JWSSignatureCreator
}

// NewAdvancedCommerceAPIInAppSignatureCreator creates a new Advanced Commerce API signature creator.
// The signingKey should be a PEM-encoded PKCS#8 ECDSA private key.
func NewAdvancedCommerceAPIInAppSignatureCreator(signingKey []byte, keyID, issuerID, bundleID string) (*AdvancedCommerceAPIInAppSignatureCreator, error) {
	base, err := NewJWSSignatureCreator("advanced-commerce-api", signingKey, keyID, issuerID, bundleID)
	if err != nil {
		return nil, err
	}
	return &AdvancedCommerceAPIInAppSignatureCreator{base}, nil
}

// CreateSignature creates a signed token for an Advanced Commerce API in-app request.
func (s *AdvancedCommerceAPIInAppSignatureCreator) CreateSignature(advancedCommerceInAppRequest any) (string, error) {
	if advancedCommerceInAppRequest == nil {
		return "", errors.New("advanced_commerce_in_app_request cannot be null")
	}

	requestJSON, err := json.Marshal(advancedCommerceInAppRequest)
	if err != nil {
		return "", err
	}

	encodedRequest := base64.StdEncoding.EncodeToString(requestJSON)
	featureSpecificClaims := map[string]any{
		"request": encodedRequest,
	}
	return s.createSignature(featureSpecificClaims)
}
