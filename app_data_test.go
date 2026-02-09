package appstore

import (
	"encoding/json"
	"testing"
)

func TestAppDataDeserialization(t *testing.T) {
	jsonData, err := readTestData("models/appData.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var appData AppData
	if err := json.Unmarshal(jsonData, &appData); err != nil {
		t.Fatalf("Failed to unmarshal AppData: %v", err)
	}

	assertEqual(t, int64(987654321), appData.AppAppleId, "AppAppleId")
	assertEqual(t, "com.example", appData.BundleId, "BundleId")
	assertEqual(t, ENVIRONMENT_SANDBOX, appData.Environment, "Environment")
	assertEqual(t, "Sandbox", appData.RawEnvironment, "RawEnvironment")
	assertEqual(t, "signed-app-transaction-info", appData.SignedAppTransactionInfo, "SignedAppTransactionInfo")
}
