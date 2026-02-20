package appstore

import (
	"testing"
)

func TestExtractTransactionIdFromAppReceipt_Empty(t *testing.T) {
	receipt, err := readTestDataString("xcode/xcode-app-receipt-empty")
	assertNoError(t, err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	assertNoError(t, err, "Failed to extract transaction ID")

	assertNil(t, transactionId, "TransactionId")
}

func TestExtractTransactionIdFromAppReceipt_WithTransaction(t *testing.T) {
	receipt, err := readTestDataString("xcode/xcode-app-receipt-with-transaction")
	assertNoError(t, err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	assertNoError(t, err, "Failed to extract transaction ID")

	assertNotNil(t, transactionId, "TransactionId")
	assertEqual(t, "0", *transactionId, "TransactionId")
}

func TestExtractTransactionIdFromTransactionReceipt(t *testing.T) {
	receipt, err := readTestDataString("mock_signed_data/legacyTransaction")
	assertNoError(t, err, "Failed to read test data")

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromTransactionReceipt(receipt)
	assertNoError(t, err, "Failed to extract transaction ID")

	assertNotNil(t, transactionId, "TransactionId")
	assertEqual(t, "33993399", *transactionId, "TransactionId")
}

func TestEncodeLength(t *testing.T) {
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
		assertEqual(t, len(test.expected), len(result), "Encoded length size")
		for i, b := range result {
			if i < len(test.expected) {
				assertEqual(t, test.expected[i], b, "Encoded length byte")
			}
		}
	}
}
