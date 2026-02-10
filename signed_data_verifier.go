package appstore

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ocsp"
)

// VerificationStatus represents the result of a verification operation.
type VerificationStatus int

const (
	OK                             VerificationStatus = 0
	VERIFICATION_FAILURE           VerificationStatus = 1
	INVALID_APP_IDENTIFIER         VerificationStatus = 2
	INVALID_CERTIFICATE            VerificationStatus = 3
	INVALID_CHAIN_LENGTH           VerificationStatus = 4
	INVALID_CHAIN                  VerificationStatus = 5
	INVALID_ENVIRONMENT            VerificationStatus = 6
	RETRYABLE_VERIFICATION_FAILURE VerificationStatus = 7
)

func (s VerificationStatus) String() string {
	switch s {
	case OK:
		return "OK"
	case VERIFICATION_FAILURE:
		return "VERIFICATION_FAILURE"
	case INVALID_APP_IDENTIFIER:
		return "INVALID_APP_IDENTIFIER"
	case INVALID_CERTIFICATE:
		return "INVALID_CERTIFICATE"
	case INVALID_CHAIN_LENGTH:
		return "INVALID_CHAIN_LENGTH"
	case INVALID_CHAIN:
		return "INVALID_CHAIN"
	case INVALID_ENVIRONMENT:
		return "INVALID_ENVIRONMENT"
	case RETRYABLE_VERIFICATION_FAILURE:
		return "RETRYABLE_VERIFICATION_FAILURE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// VerificationException is an error that indicates verification failed with a specific status.
type VerificationException struct {
	Status VerificationStatus
	Err    error
}

func (e *VerificationException) Error() string {
	return fmt.Sprintf("Verification failed with status %s: %v", e.Status, e.Err)
}

// NewVerificationException creates a new verification exception with the given status and error.
func NewVerificationException(status VerificationStatus, err error) *VerificationException {
	return &VerificationException{Status: status, Err: err}
}

// SignedDataVerifier provides utility methods for verifying and decoding App Store signed data.
type SignedDataVerifier struct {
	chainVerifier       *chainVerifier
	environment         Environment
	bundleID            string
	appAppleID          int64
	enableOnlineChecks  bool
	allowAnyEnvironment bool
}

// NewSignedDataVerifier creates a new SignedDataVerifier for verifying App Store signed data.
// The appAppleID is required when the environment is Production.
//
// See https://developer.apple.com/documentation/appstoreserverapi
func NewSignedDataVerifier(rootCertificates [][]byte, enableOnlineChecks bool, environment Environment, bundleID string, appAppleID int64) (*SignedDataVerifier, error) {
	if environment == ENVIRONMENT_PRODUCTION && appAppleID == 0 {
		return nil, errors.New("appAppleId is required when the environment is Production")
	}

	cv, err := newChainVerifier(rootCertificates)
	if err != nil {
		return nil, err
	}
	return &SignedDataVerifier{
		chainVerifier:      cv,
		environment:        environment,
		bundleID:           bundleID,
		appAppleID:         appAppleID,
		enableOnlineChecks: enableOnlineChecks,
	}, nil
}

// VerifyAndDecodeRenewalInfo verifies and decodes a signedRenewalInfo obtained from the App Store Server API,
// an App Store Server Notification, or from a device.
//
// See https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfo
func (v *SignedDataVerifier) VerifyAndDecodeRenewalInfo(signedRenewalInfo string) (*JWSRenewalInfoDecodedPayload, error) {
	payload := &JWSRenewalInfoDecodedPayload{}
	if err := v.decodeSignedObject(signedRenewalInfo, payload); err != nil {
		return nil, err
	}
	if payload.Environment != v.environment && !v.allowAnyEnvironment {
		return nil, NewVerificationException(INVALID_ENVIRONMENT, errors.New("environment mismatch"))
	}
	return payload, nil
}

// VerifyAndDecodeSignedTransaction verifies and decodes a signedTransaction obtained from the App Store Server API,
// an App Store Server Notification, or from a device.
//
// See https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
func (v *SignedDataVerifier) VerifyAndDecodeSignedTransaction(signedTransaction string) (*JWSTransactionDecodedPayload, error) {
	payload := &JWSTransactionDecodedPayload{}
	if err := v.decodeSignedObject(signedTransaction, payload); err != nil {
		return nil, err
	}
	if payload.BundleId != v.bundleID {
		return nil, NewVerificationException(INVALID_APP_IDENTIFIER, errors.New("bundleId mismatch"))
	}
	if payload.Environment != v.environment && !v.allowAnyEnvironment {
		return nil, NewVerificationException(INVALID_ENVIRONMENT, errors.New("environment mismatch"))
	}
	return payload, nil
}

// VerifyAndDecodeNotification verifies and decodes an App Store Server Notification signedPayload.
//
// See https://developer.apple.com/documentation/appstoreservernotifications/signedpayload
func (v *SignedDataVerifier) VerifyAndDecodeNotification(signedPayload string) (*ResponseBodyV2DecodedPayload, error) {
	payload := &ResponseBodyV2DecodedPayload{}
	if err := v.decodeSignedObject(signedPayload, payload); err != nil {
		return nil, err
	}

	var bundleID string
	var appAppleID int64
	var environment Environment

	switch {
	case payload.Data != nil:
		bundleID = payload.Data.BundleId
		appAppleID = payload.Data.AppAppleId
		environment = payload.Data.Environment
	case payload.Summary != nil:
		bundleID = payload.Summary.BundleId
		appAppleID = payload.Summary.AppAppleId
		environment = payload.Summary.Environment
	case payload.ExternalPurchaseToken != nil:
		bundleID = payload.ExternalPurchaseToken.BundleId
		appAppleID = payload.ExternalPurchaseToken.AppAppleId
		if strings.HasPrefix(payload.ExternalPurchaseToken.ExternalPurchaseId, "SANDBOX") {
			environment = ENVIRONMENT_SANDBOX
		} else {
			environment = ENVIRONMENT_PRODUCTION
		}
	case payload.AppData != nil:
		bundleID = payload.AppData.BundleId
		appAppleID = payload.AppData.AppAppleId
		environment = payload.AppData.Environment
	}

	if err := v.verifyNotification(bundleID, appAppleID, environment); err != nil {
		return nil, err
	}
	return payload, nil
}

func (v *SignedDataVerifier) verifyNotification(bundleID string, appAppleID int64, environment Environment) error {
	if bundleID != v.bundleID || (v.environment == ENVIRONMENT_PRODUCTION && appAppleID != v.appAppleID) {
		return NewVerificationException(INVALID_APP_IDENTIFIER, errors.New("app identifier mismatch"))
	}
	if environment != v.environment && !v.allowAnyEnvironment {
		return NewVerificationException(INVALID_ENVIRONMENT, errors.New("environment mismatch"))
	}
	return nil
}

// VerifyAndDecodeAppTransaction verifies and decodes a signed AppTransaction.
//
// See https://developer.apple.com/documentation/storekit/apptransaction
func (v *SignedDataVerifier) VerifyAndDecodeAppTransaction(signedAppTransaction string) (*AppTransaction, error) {
	payload := &AppTransaction{}
	if err := v.decodeSignedObject(signedAppTransaction, payload); err != nil {
		return nil, err
	}
	environment := payload.ReceiptType
	if payload.BundleId != v.bundleID || (v.environment == ENVIRONMENT_PRODUCTION && payload.AppAppleId != v.appAppleID) {
		return nil, NewVerificationException(INVALID_APP_IDENTIFIER, errors.New("app identifier mismatch"))
	}
	if environment != v.environment {
		return nil, NewVerificationException(INVALID_ENVIRONMENT, errors.New("environment mismatch"))
	}
	return payload, nil
}

// VerifyAndDecodeRealtimeRequest verifies and decodes a Retention Messaging API signedPayload.
//
// See https://developer.apple.com/documentation/retentionmessaging/signedpayload
func (v *SignedDataVerifier) VerifyAndDecodeRealtimeRequest(signedPayload string) (*DecodedRealtimeRequestBody, error) {
	payload := &DecodedRealtimeRequestBody{}
	if err := v.decodeSignedObject(signedPayload, payload); err != nil {
		return nil, err
	}
	if v.environment == ENVIRONMENT_PRODUCTION && payload.AppAppleId != v.appAppleID {
		return nil, NewVerificationException(INVALID_APP_IDENTIFIER, errors.New("app identifier mismatch"))
	}
	if payload.Environment != v.environment {
		return nil, NewVerificationException(INVALID_ENVIRONMENT, errors.New("environment mismatch"))
	}
	return payload, nil
}

func (v *SignedDataVerifier) decodeSignedObject(signedObj string, destination any) error {
	claims := jwt.MapClaims{}
	var err error
	if v.environment == ENVIRONMENT_XCODE || v.environment == ENVIRONMENT_LOCAL_TESTING {
		_, _, err = new(jwt.Parser).ParseUnverified(signedObj, &claims)
	} else {
		_, err = jwt.ParseWithClaims(signedObj, &claims, func(token *jwt.Token) (any, error) {
			x5c, ok := token.Header["x5c"].([]any)
			if !ok || len(x5c) == 0 {
				return nil, errors.New("x5c header is missing or empty")
			}

			certs := make([]string, len(x5c))
			for i, v := range x5c {
				certs[i], ok = v.(string)
				if !ok {
					return nil, errors.New("invalid x5c header format")
				}
			}

			alg, ok := token.Header["alg"].(string)
			if !ok || alg != "ES256" {
				return nil, errors.New("invalid algorithm header")
			}

			var signedDateMillis int64
			if sd, ok := claims["signedDate"].(int64); ok {
				signedDateMillis = sd
			}
			if signedDateMillis == 0 {
				if rcd, ok := claims["receiptCreationDate"].(int64); ok {
					signedDateMillis = rcd
				}
			}

			var effectiveDate time.Time
			if v.enableOnlineChecks || signedDateMillis == 0 {
				effectiveDate = time.Now()
			} else {
				effectiveDate = time.UnixMilli(signedDateMillis)
			}

			publicKey, err := v.chainVerifier.verifyChain(certs, v.enableOnlineChecks, effectiveDate)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		}, jwt.WithValidMethods([]string{"ES256"}))
	}

	if err != nil {
		return NewVerificationException(VERIFICATION_FAILURE, err)
	}

	data, err := json.Marshal(claims)
	if err != nil {
		return NewVerificationException(VERIFICATION_FAILURE, err)
	}
	return json.Unmarshal(data, destination)
}

type cacheEntry struct {
	publicKey *ecdsa.PublicKey
	expiry    time.Time
}

type chainVerifier struct {
	rootCertificates *x509.CertPool
	cache            map[string]cacheEntry
	cacheMutex       sync.RWMutex
}

const (
	maxCacheSize   = 32
	cacheTimeLimit = 15 * time.Minute
)

func newChainVerifier(rootCerts [][]byte) (*chainVerifier, error) {
	pool := x509.NewCertPool()
	for _, certBytes := range rootCerts {
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, NewVerificationException(INVALID_CERTIFICATE, err)
		}
		pool.AddCert(cert)
	}
	return &chainVerifier{
		rootCertificates: pool,
		cache:            make(map[string]cacheEntry),
	}, nil
}

func (cv *chainVerifier) verifyChain(certificates []string, performOnlineChecks bool, effectiveDate time.Time) (*ecdsa.PublicKey, error) {
	if performOnlineChecks {
		cacheKey := strings.Join(certificates, "|")
		cv.cacheMutex.RLock()
		if entry, ok := cv.cache[cacheKey]; ok && time.Now().Before(entry.expiry) {
			cv.cacheMutex.RUnlock()
			return entry.publicKey, nil
		}
		cv.cacheMutex.RUnlock()
	}

	if len(certificates) != 3 {
		return nil, NewVerificationException(INVALID_CHAIN_LENGTH, errors.New("invalid chain length"))
	}

	var parsedCerts []*x509.Certificate
	for _, certStr := range certificates {
		certBytes, err := base64.StdEncoding.DecodeString(certStr)
		if err != nil {
			return nil, NewVerificationException(INVALID_CERTIFICATE, err)
		}
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, NewVerificationException(INVALID_CERTIFICATE, err)
		}
		parsedCerts = append(parsedCerts, cert)
	}

	leaf := parsedCerts[0]
	intermediates := x509.NewCertPool()
	for _, cert := range parsedCerts[1:] {
		intermediates.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots:         cv.rootCertificates,
		Intermediates: intermediates,
		CurrentTime:   effectiveDate,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	if _, err := leaf.Verify(opts); err != nil {
		return nil, NewVerificationException(VERIFICATION_FAILURE, err)
	}

	// OID checks
	if err := cv.checkOID(leaf, "1.2.840.113635.100.6.11.1"); err != nil {
		return nil, err
	}
	if err := cv.checkOID(parsedCerts[1], "1.2.840.113635.100.6.2.1"); err != nil {
		return nil, err
	}

	// OCSP check
	if performOnlineChecks {
		// Retrieve the verified chain [leaf, intermediate, root] for OCSP validation
		chains, err := leaf.Verify(opts)
		if err != nil || len(chains) == 0 || len(chains[0]) < 3 {
			return nil, NewVerificationException(VERIFICATION_FAILURE, errors.New("failed to verify chain for OCSP"))
		}
		verifiedChain := chains[0]
		if err := cv.checkOCSP(verifiedChain[1], verifiedChain[2], verifiedChain[2]); err != nil {
			return nil, err
		}
		if err := cv.checkOCSP(verifiedChain[0], verifiedChain[1], verifiedChain[2]); err != nil {
			return nil, err
		}

		// Update cache
		pubKey, ok := leaf.PublicKey.(*ecdsa.PublicKey)
		if ok {
			cacheKey := strings.Join(certificates, "|")
			cv.saveToCache(cacheKey, pubKey)
		}
	}

	pubKey, ok := leaf.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, NewVerificationException(VERIFICATION_FAILURE, errors.New("not an ECDSA public key"))
	}

	return pubKey, nil
}

func (cv *chainVerifier) saveToCache(cacheKey string, pubKey *ecdsa.PublicKey) {
	cv.cacheMutex.Lock()
	defer cv.cacheMutex.Unlock()

	// Basic eviction if full
	if len(cv.cache) >= maxCacheSize {
		for k, v := range cv.cache {
			if time.Now().After(v.expiry) {
				delete(cv.cache, k)
			}
		}
		// If still full, just delete a random one
		if len(cv.cache) >= maxCacheSize {
			for k := range cv.cache {
				delete(cv.cache, k)
				break
			}
		}
	}
	cv.cache[cacheKey] = cacheEntry{
		publicKey: pubKey,
		expiry:    time.Now().Add(cacheTimeLimit),
	}
}

func (cv *chainVerifier) checkOID(cert *x509.Certificate, expectedOID string) error {
	for _, ext := range cert.Extensions {
		if ext.Id.String() == expectedOID {
			return nil
		}
	}
	return NewVerificationException(VERIFICATION_FAILURE, fmt.Errorf("missing expected OID: %s", expectedOID))
}

func (cv *chainVerifier) checkOCSP(cert, issuer, root *x509.Certificate) error {
	for _, server := range cert.OCSPServer {
		if err := cv.checkOCSPServer(server, cert, issuer, root); err == nil {
			return nil
		}
	}
	return NewVerificationException(VERIFICATION_FAILURE, errors.New("failed to get a valid OCSP response"))
}

func (cv *chainVerifier) checkOCSPServer(server string, cert, issuer, root *x509.Certificate) error {
	opts := &ocsp.RequestOptions{Hash: crypto.SHA256}
	buffer, err := ocsp.CreateRequest(cert, issuer, opts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", server, bytes.NewReader(buffer))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/ocsp-request")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return NewVerificationException(RETRYABLE_VERIFICATION_FAILURE, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("OCSP server returned non-200 status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewVerificationException(RETRYABLE_VERIFICATION_FAILURE, err)
	}

	ocspResp, err := ocsp.ParseResponse(body, nil)
	if err != nil {
		return err
	}

	if ocspResp.Status != ocsp.Good {
		return errors.New("OCSP status is not Good")
	}

	// Check serial number
	if ocspResp.SerialNumber.Cmp(cert.SerialNumber) != 0 {
		return errors.New("OCSP serial number mismatch")
	}

	if err := cv.verifyOCSPResponseSignature(ocspResp, issuer, root); err != nil {
		return err
	}

	return nil
}

func (cv *chainVerifier) verifyOCSPResponseSignature(ocspResp *ocsp.Response, issuer, root *x509.Certificate) error {
	// If response doesn't include a certificate, it must be signed by the issuer
	if ocspResp.Certificate == nil {
		return ocspResp.CheckSignatureFrom(issuer)
	}

	// Response includes a certificate; verified chain includes leaf (responder), intermediate (issuer), and root.
	// ocsp.ParseResponse has already identified the responder cert.
	roots := x509.NewCertPool()
	roots.AddCert(root)
	intermediates := x509.NewCertPool()
	intermediates.AddCert(issuer)

	_, err := ocspResp.Certificate.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageOCSPSigning}, // Validates OCSP signing usage
	})
	if err != nil {
		return err
	}

	return ocspResp.CheckSignatureFrom(ocspResp.Certificate)
}
