package appstore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TEST_BUNDLE_ID = "com.example"
	TEST_KEY_ID    = "ABCDEFG123"
	TEST_ISSUER_ID = "01234567-890a-bcde-f012-345678901234"
)

// readTestData reads a test data file from testdata directory
func readTestData(relativePath string) ([]byte, error) {
	fullPath := filepath.Join("testdata", relativePath)
	return os.ReadFile(fullPath)
}

// readTestDataString reads a test data file as a string
func readTestDataString(relativePath string) (string, error) {
	data, err := readTestData(relativePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// createTestSignedDataVerifier creates a SignedDataVerifier for testing
// with the test CA certificate and no OCSP validation
func createTestSignedDataVerifier(env Environment, bundleID string, appAppleID *int64) (*SignedDataVerifier, error) {
	testCA, err := readTestData("certs/testCA.der")
	if err != nil {
		return nil, err
	}

	var appID int64
	if appAppleID != nil {
		appID = *appAppleID
	}

	return NewSignedDataVerifier([][]byte{testCA}, false, env, bundleID, appID)
}

// createSignedDataFromJSON creates a signed JWT token from JSON test data
// This generates a self-signed token for testing purposes
func createSignedDataFromJSON(jsonPath string) (string, error) {
	// Read the JSON payload
	jsonData, err := readTestData(jsonPath)
	if err != nil {
		return "", err
	}

	// Parse JSON into a map
	var payload map[string]any
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return "", err
	}

	// Generate a temporary EC private key for signing
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(payload))

	// Sign the token
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// decodeJWTWithoutVerification decodes a JWT token without signature verification
// Returns the header and payload as maps
func decodeJWTWithoutVerification(tokenString string) (header, payload map[string]any, err error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, nil, err
	}

	// Get header
	if token.Header != nil {
		header = token.Header
	}

	// Get payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		payload = claims
	}

	return header, payload, nil
}

// int64Ptr returns a pointer to an int64 value
func int64Ptr(v int64) *int64 {
	return &v
}

// stringPtr returns a pointer to a string value
func stringPtr(v string) *string {
	return &v
}

// assertEqual is a helper function for assertions in tests
func assertEqual(t *testing.T, expected, actual any, fieldName string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", fieldName, expected, actual)
	}
}

// assertNoError is a helper function to assert that an error is nil
func assertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", message, err)
	}
}

// assertError is a helper function to assert that an error is not nil
func assertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error, got nil", message)
	}
}

// assertTrue is a helper function to assert that a value is true
func assertTrue(t *testing.T, actual bool, fieldName string) {
	t.Helper()
	if !actual {
		t.Errorf("%s: expected true, got false", fieldName)
	}
}

// assertNil is a helper function to assert that a value is nil
func assertNil(t *testing.T, actual any, fieldName string) {
	t.Helper()
	if !isNil(actual) {
		t.Errorf("%s: expected nil, got %v", fieldName, actual)
	}
}

// assertNotNil is a helper function to assert that a value is not nil
func assertNotNil(t *testing.T, actual any, fieldName string) {
	t.Helper()
	if isNil(actual) {
		t.Fatalf("%s: expected not nil, got nil", fieldName)
	}
}

func isNil(i any) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}
