# App Store Server Library (Go)

[![CI](https://github.com/laishere/app-store-server-library-go/actions/workflows/ci.yml/badge.svg)](https://github.com/laishere/app-store-server-library-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/laishere/app-store-server-library-go/branch/main/graph/badge.svg)](https://codecov.io/gh/laishere/app-store-server-library-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/laishere/app-store-server-library-go)](https://goreportcard.com/report/github.com/laishere/app-store-server-library-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/laishere/app-store-server-library-go.svg)](https://pkg.go.dev/github.com/laishere/app-store-server-library-go)
[![GitHub release](https://img.shields.io/github/v/release/laishere/app-store-server-library-go)](https://github.com/laishere/app-store-server-library-go/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/laishere/app-store-server-library-go)](https://github.com/laishere/app-store-server-library-go/blob/main/go.mod)

A Go port of the [Apple App Store Server Library](https://github.com/apple/app-store-server-library-python).

This library provides a Go implementation for interacting with the [App Store Server API](https://developer.apple.com/documentation/appstoreserverapi), [App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications), and [Retention Messaging API](https://developer.apple.com/documentation/retentionmessaging).

## Table of Contents
1. [Installation](#installation)
2. [Documentation](#documentation)
3. [Usage](#usage)
4. [License](#license)

## Installation

```bash
go get github.com/laishere/app-store-server-library-go
```

## Documentation

[WWDC Video](https://developer.apple.com/videos/play/wwdc2023/10143/)

### Obtaining an In-App Purchase key from App Store Connect

To use the App Store Server API or create promotional offer signatures, a signing key downloaded from App Store Connect is required. To obtain this key, you must have the Admin role. Go to Users and Access > Integrations > In-App Purchase. Here you can create and manage keys, as well as find your issuer ID. When using a key, you'll need the key ID and issuer ID as well.

### Obtaining Apple Root Certificates

Download and store the root certificates found in the Apple Root Certificates section of the [Apple PKI](https://www.apple.com/certificateauthority/) site. Provide these certificates as an array to a `SignedDataVerifier` to allow verifying the signed data comes from Apple.

## Usage

For more detailed examples, see the [examples](examples) directory.

### API Client

```go
signingKey, _ := os.ReadFile("SubscriptionKey_ABCD123456.p8")
client, _ := appstore.NewAPIClient(signingKey, "ABCD123456", "issuer_id", "com.example", appstore.ENVIRONMENT_SANDBOX)

response, err := client.RequestTestNotification()
```

### Verification Usage

```go
rootCert, _ := os.ReadFile("AppleRootCA-G3.cer")
verifier, _ := appstore.NewSignedDataVerifier([][]byte{rootCert}, true, appstore.ENVIRONMENT_SANDBOX, "com.example", 123456789)

payload, err := verifier.VerifyAndDecodeNotification(signedPayload)
```

### Receipt Usage

```go
utility := appstore.NewReceiptUtility()
transactionId, err := utility.ExtractTransactionIdFromAppReceipt(appReceipt)
```

### Promotional Offer Signature Creation

```go
creator, _ := appstore.NewPromotionalOfferV2SignatureCreator(signingKey, "keyId", "issuerId", "bundleId")
signature, err := creator.CreateSignature(productID, offerIdentifier, nil)
```

## Features

- **App Store Server API Client**: Complete implementation of the App Store Server API endpoints
  - Transaction history and subscription status
  - Order lookup and refund history  
  - Test notifications
  - Consumption information
  - Subscription renewal date extensions
  - Notification history
- **App Store Server Notifications**: Verify and decode App Store Server Notifications V2
- **Retention Messaging API**: Upload and manage retention messaging images and messages
- **Receipt Utility**: Extract transaction IDs from App Receipts and transactional receipts
- **Signed Data Verification**: Verify and decode JWS signed data from the App Store
- **Signature Creators**: Generate signatures for various use cases
  - Promotional Offer V2 signatures
  - Introductory offer eligibility signatures
  - Advanced Commerce API in-app signatures

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
Portions copyright (c) 2023 Apple Inc.
