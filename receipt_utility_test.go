package appstore

import (
	"testing"
)

func TestExtractTransactionIdFromAppReceipt_Empty(t *testing.T) {
	receipt, err := readTestDataString("xcode/xcode-app-receipt-empty")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to extract transaction ID: %v", err)
	}

	if transactionId != nil {
		t.Errorf("Expected nil transaction ID for empty receipt, got %s", *transactionId)
	}
}

func TestExtractTransactionIdFromAppReceipt_WithTransaction(t *testing.T) {
	receipt, err := readTestDataString("xcode/xcode-app-receipt-with-transaction")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromAppReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to extract transaction ID: %v", err)
	}

	if transactionId == nil {
		t.Fatal("Expected transaction ID, got nil")
	}
	assertEqual(t, "0", *transactionId, "TransactionId")
}

func TestExtractTransactionIdFromTransactionReceipt(t *testing.T) {
	receipt, err := readTestDataString("mock_signed_data/legacyTransaction")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	utility := NewReceiptUtility()
	transactionId, err := utility.ExtractTransactionIdFromTransactionReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to extract transaction ID: %v", err)
	}

	if transactionId == nil {
		t.Fatal("Expected transaction ID, got nil")
	}
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
		if len(result) != len(test.expected) {
			t.Errorf("For length %d, expected result length %d, got %d", test.length, len(test.expected), len(result))
			continue
		}
		for i, b := range result {
			if b != test.expected[i] {
				t.Errorf("For length %d, at index %d expected byte %x, got %x", test.length, i, test.expected[i], b)
				break
			}
		}
	}
}
