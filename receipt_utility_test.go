package appstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractTransactionIdFromAppReceipt_Empty(t *testing.T) {
	assert := assert.New(t)
	receipt, err := readTestDataString("xcode/xcode-app-receipt-empty")
	assert.NoError(err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	assert.NoError(err, "Failed to extract transaction ID")

	assert.Nil(transactionId, "TransactionId")
}

func TestExtractTransactionIdFromAppReceipt_WithTransaction(t *testing.T) {
	assert := assert.New(t)
	receipt, err := readTestDataString("xcode/xcode-app-receipt-with-transaction")
	assert.NoError(err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	assert.NoError(err, "Failed to extract transaction ID")

	assert.NotNil(transactionId, "TransactionId")
	assert.Equal("0", *transactionId, "TransactionId")
}

func TestExtractTransactionIdFromTransactionReceipt(t *testing.T) {
	assert := assert.New(t)
	receipt, err := readTestDataString("mock_signed_data/legacyTransaction")
	assert.NoError(err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromTransactionReceipt(receipt)
	assert.NoError(err, "Failed to extract transaction ID")

	assert.NotNil(transactionId, "TransactionId")
	assert.Equal("33993399", *transactionId, "TransactionId")
}

func TestEncodeLength(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		length   int
		expected []byte
	}{
		{0, []byte{0}},
		{127, []byte{127}},
		{128, []byte{0x81, 0x80}},
		{255, []byte{0x81, 0xFF}},
		{256, []byte{0x82, 0x01, 0x00}},
		{65535, []byte{0x82, 0xFF, 0xFF}},
		{65536, []byte{0x83, 0x01, 0x00, 0x00}},
	}

	for _, test := range tests {
		result := encodeLength(test.length)
		assert.Equal(len(test.expected), len(result), "Encoded length size")
		for i, b := range result {
			if i < len(test.expected) {
				assert.Equal(test.expected[i], b, "Encoded length byte")
			}
		}
	}
}
