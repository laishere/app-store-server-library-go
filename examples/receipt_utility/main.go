package main

import (
	"fmt"
	"os"

	appstore "github.com/laishere/app-store-server-library-go"
)

func main() {
	// 1. Setup API Client (optional, but shows how to use the extracted ID)
	signingKey, _ := os.ReadFile("SubscriptionKey_ABCD123456.p8")
	client, _ := appstore.NewAPIClient(signingKey, "keyId", "issuerId", "bundleId", appstore.ENVIRONMENT_SANDBOX)

	// 2. Setup Receipt Utility
	utility := appstore.NewReceiptUtility()

	// Base64 encoded string from the app's 'appStoreReceiptURL'
	appReceipt := "MIITYQYJKoZIhvcNAQcCoIITUjCCE1ACAQExCzAJBgUrDgMCGgUAMIIB..."

	// 3. Extract the transaction ID
	// This works for both StoreKit 1 receipts and Xcode local receipts
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(appReceipt)
	if err != nil {
		fmt.Printf("Extraction error: %v\n", err)
		return
	}

	if transactionId == nil {
		fmt.Println("No transaction found in the receipt.")
		return
	}

	fmt.Printf("Extracted Transaction ID: %s\n", *transactionId)

	// 4. Use the ID to query the App Store Server API
	// This is a common pattern: get ID from device receipt -> query server API for latest status
	response, err := client.GetTransactionHistory(*transactionId, nil, "", appstore.GET_TRANSACTION_HISTORY_VERSION_V2)
	if err != nil {
		fmt.Printf("API Query failed: %v\n", err)
		return
	}

	fmt.Printf("Successfully retrieved %d signed transactions from history\n", len(response.SignedTransactions))
}
