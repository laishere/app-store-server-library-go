package main

import (
	"fmt"
	"net/url"
	"os"

	appstore "github.com/laishere/app-store-server-library-go"
)

func main() {
	signingKey, err := os.ReadFile("SubscriptionKey_ABCD123456.p8")
	if err != nil {
		fmt.Printf("Failed to read key: %v\n", err)
		return
	}

	keyId := "ABCD123456"
	issuerId := "57246542-96fe-1a63-e053-0824d011072a"
	bundleId := "com.example.app"

	client, err := appstore.NewAPIClient(
		signingKey,
		keyId,
		issuerId,
		bundleId,
		appstore.ENVIRONMENT_SANDBOX,
	)
	if err != nil {
		panic(err)
	}

	// Example: Fetching FULL Transaction History with Pagination
	transactionID := "original_transaction_id"
	revision := ""
	var allTransactions []string

	for {
		// Use url.Values to filter by product type if needed
		queryParams := url.Values{}
		queryParams.Set("productType", "AUTO_RENEWABLE")

		response, err := client.GetTransactionHistory(
			transactionID,
			queryParams,
			revision,
			appstore.GET_TRANSACTION_HISTORY_VERSION_V2,
		)
		if err != nil {
			fmt.Printf("Error during fetch: %v\n", err)
			break
		}

		allTransactions = append(allTransactions, response.SignedTransactions...)

		if !response.HasMore {
			break
		}
		// Update revision to fetch the next page
		revision = response.Revision
		fmt.Printf("Fetched a page, total so far: %d\n", len(allTransactions))
	}

	fmt.Printf("Final count of transactions: %d\n", len(allTransactions))
}
