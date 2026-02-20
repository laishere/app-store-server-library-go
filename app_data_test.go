package appstore

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppDataDeserialization(t *testing.T) {
	assert := assert.New(t)
	jsonData, err := readTestData("models/appData.json")
	assert.NoError(err, "Failed to read test data")

	var appData AppData
	err = json.Unmarshal(jsonData, &appData)
	assert.NoError(err, "Failed to unmarshal AppData")

	assert.Equal(int64(987654321), appData.AppAppleId, "AppAppleId")
	assert.Equal("com.example", appData.BundleId, "BundleId")
	assert.Equal(ENVIRONMENT_SANDBOX, appData.Environment, "Environment")
	assert.Equal("signed-app-transaction-info", appData.SignedAppTransactionInfo, "SignedAppTransactionInfo")
}
