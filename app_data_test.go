package appstore

import (
	"encoding/json"
	"testing"
)

func TestAppDataDeserialization(t *testing.T) {
	jsonData, err := readTestData("models/appData.json")
	assertNoError(t, err, "Failed to read test data")

	var appData AppData
	err = json.Unmarshal(jsonData, &appData)
	assertNoError(t, err, "Failed to unmarshal AppData")

	assertEqual(t, int64(987654321), appData.AppAppleId, "AppAppleId")
	assertEqual(t, "com.example", appData.BundleId, "BundleId")
	assertEqual(t, ENVIRONMENT_SANDBOX, appData.Environment, "Environment")
	assertEqual(t, "signed-app-transaction-info", appData.SignedAppTransactionInfo, "SignedAppTransactionInfo")
}
