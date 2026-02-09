package main

import (
	"fmt"
	"os"

	appstore "github.com/laishere/app-store-server-library-go"
)

func main() {
	// 1. Load the Apple Root Certificate.
	// You should download 'AppleRootCA-G3.cer' from Apple PKI (https://www.apple.com/certificateauthority/)
	rootCert, err := os.ReadFile("AppleRootCA-G3.cer")
	if err != nil {
		fmt.Printf("Failed to read root certificate: %v\n", err)
		fmt.Println("Please download 'AppleRootCA-G3.cer' from Apple PKI website.")
		return
	}

	// 2. Initialize the Verifier
	// ONLINE checks (OCSP) are enabled by default for production systems.
	verifier, err := appstore.NewSignedDataVerifier(
		[][]byte{rootCert},
		true, // Enable online checks
		appstore.ENVIRONMENT_SANDBOX,
		"com.example.app",
		123456789, // App Apple ID (find this in App Store Connect)
	)
	if err != nil {
		panic(err)
	}

	// 3. Receive the signed payload (e.g., from a Webhook)
	signedPayload := "eyJhbGciOiJFUzI1NiIsIng1YyI6WyJNSUlCTURDQ0FSUm... (the raw JWS string)"

	// 4. Verify and Decode
	payload, err := verifier.VerifyAndDecodeNotification(signedPayload)
	if err != nil {
		if verifErr, ok := err.(*appstore.VerificationException); ok {
			fmt.Printf("Verification specialized error: %v (Status: %s)\n", verifErr.Err, verifErr.Status)
		} else {
			fmt.Printf("General error: %v\n", err)
		}
		return
	}

	fmt.Printf("Verified Notification: %s\n", payload.NotificationType)
	fmt.Printf("Signed Date: %d\n", payload.SignedDate)
	if payload.Data != nil {
		fmt.Printf("Bundle ID: %s\n", payload.Data.BundleId)
	}
}
