package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	ROOT_CA_BASE64_ENCODED                                   = "MIIBgjCCASmgAwIBAgIJALUc5ALiH5pbMAoGCCqGSM49BAMDMDYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDdXBlcnRpbm8wHhcNMjMwMTA1MjEzMDIyWhcNMzMwMTAyMjEzMDIyWjA2MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJQ3VwZXJ0aW5vMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEc+/Bl+gospo6tf9Z7io5tdKdrlN1YdVnqEhEDXDShzdAJPQijamXIMHf8xWWTa1zgoYTxOKpbuJtDplz1XriTaMgMB4wDAYDVR0TBAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwCgYIKoZIzj0EAwMDRwAwRAIgemWQXnMAdTad2JDJWng9U4uBBL5mA7WI05H7oH7c6iQCIHiRqMjNfzUAyiu9h6rOU/K+iTR0I/3Y/NSWsXHX+acc"
	INTERMEDIATE_CA_BASE64_ENCODED                           = "MIIBnzCCAUWgAwIBAgIBCzAKBggqhkjOPQQDAzA2MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJQ3VwZXJ0aW5vMB4XDTIzMDEwNTIxMzEwNVoXDTMzMDEwMTIxMzEwNVowRTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlDdXBlcnRpbm8xFTATBgNVBAoMDEludGVybWVkaWF0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABBUN5V9rKjfRiMAIojEA0Av5Mp0oF+O0cL4gzrTF178inUHugj7Et46NrkQ7hKgMVnjogq45Q1rMs+cMHVNILWqjNTAzMA8GA1UdEwQIMAYBAf8CAQAwDgYDVR0PAQH/BAQDAgEGMBAGCiqGSIb3Y2QGAgEEAgUAMAoGCCqGSM49BAMDA0gAMEUCIQCmsIKYs41ullssHX4rVveUT0Z7Is5/hLK1lFPTtun3hAIgc2+2RG5+gNcFVcs+XJeEl4GZ+ojl3ROOmll+ye7dynQ="
	LEAF_CERT_BASE64_ENCODED                                 = "MIIBoDCCAUagAwIBAgIBDDAKBggqhkjOPQQDAzBFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCUN1cGVydGlubzEVMBMGA1UECgwMSW50ZXJtZWRpYXRlMB4XDTIzMDEwNTIxMzEzNFoXDTMzMDEwMTIxMzEzNFowPTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlDdXBlcnRpbm8xDTALBgNVBAoMBExlYWYwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATitYHEaYVuc8g9AjTOwErMvGyPykPa+puvTI8hJTHZZDLGas2qX1+ErxgQTJgVXv76nmLhhRJH+j25AiAI8iGsoy8wLTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDAQBgoqhkiG92NkBgsBBAIFADAKBggqhkjOPQQDAwNIADBFAiBX4c+T0Fp5nJ5QRClRfu5PSByRvNPtuaTsk0vPB3WAIAIhANgaauAj/YP9s0AkEhyJhxQO/6Q2zouZ+H1CIOehnMzQ"
	INTERMEDIATE_CA_INVALID_OID_BASE64_ENCODED               = "MIIBnjCCAUWgAwIBAgIBDTAKBggqhkjOPQQDAzA2MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJQ3VwZXJ0aW5vMB4XDTIzMDEwNTIxMzYxNFoXDTMzMDEwMTIxMzYxNFowRTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlDdXBlcnRpbm8xFTATBgNVBAoMDEludGVybWVkaWF0ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABBUN5V9rKjfRiMAIojEA0Av5Mp0oF+O0cL4gzrTF178inUHugj7Et46NrkQ7hKgMVnjogq45Q1rMs+cMHVNILWqjNTAzMA8GA1UdEwQIMAYBAf8CAQAwDgYDVR0PAQH/BAQDAgEGMBAGCiqGSIb3Y2QGAgIEAgUAMAoGCCqGSM49BAMDA0cAMEQCIFROtTE+RQpKxNXETFsf7Mc0h+5IAsxxo/X6oCC/c33qAiAmC5rn5yCOOEjTY4R1H1QcQVh+eUwCl13NbQxWCuwxxA=="
	LEAF_CERT_FOR_INTERMEDIATE_CA_INVALID_OID_BASE64_ENCODED = "MIIBnzCCAUagAwIBAgIBDjAKBggqhkjOPQQDAzBFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCUN1cGVydGlubzEVMBMGA1UECgwMSW50ZXJtZWRpYXRlMB4XDTIzMDEwNTIxMzY1OFoXDTMzMDEwMTIxMzY1OFowPTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlDdXBlcnRpbm8xDTALBgNVBAoMBExlYWYwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATitYHEaYVuc8g9AjTOwErMvGyPykPa+puvTI8hJTHZZDLGas2qX1+ErxgQTJgVXv76nmLhhRJH+j25AiAI8iGsoy8wLTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDAQBgoqhkiG92NkBgsBBAIFADAKBggqhkjOPQQDAwNHADBEAiAUAs+gzYOsEXDwQquvHYbcVymyNqDtGw9BnUFp2YLuuAIgXxQ3Ie9YU0cMqkeaFd+lyo0asv9eyzk6stwjeIeOtTU="
	LEAF_CERT_INVALID_OID_BASE64_ENCODED                     = "MIIBoDCCAUagAwIBAgIBDzAKBggqhkjOPQQDAzBFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCUN1cGVydGlubzEVMBMGA1UECgwMSW50ZXJtZWRpYXRlMB4XDTIzMDEwNTIxMzczMVoXDTMzMDEwMTIxMzczMVowPTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlDdXBlcnRpbm8xDTALBgNVBAoMBExlYWYwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATitYHEaYVuc8g9AjTOwErMvGyPykPa+puvTI8hJTHZZDLGas2qX1+ErxgQTJgVXv76nmLhhRJH+j25AiAI8iGsoy8wLTAJBgNVHRMEAjAAMA4GA1UdDwEB/wQEAwIHgDAQBgoqhkiG92NkBgsCBAIFADAKBggqhkjOPQQDAwNIADBFAiAb+7S3i//bSGy7skJY9+D4VgcQLKFeYfIMSrUCmdrFqwIhAIMVwzD1RrxPRtJyiOCXLyibIvwcY+VS73HYfk0O9lgz"
	REAL_APPLE_ROOT_BASE64_ENCODED                           = "MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwSQXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcNMTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBSb290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtfTjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySrMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gAMGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM6BgD56KyKA=="
	REAL_APPLE_INTERMEDIATE_BASE64_ENCODED                   = "MIIDFjCCApygAwIBAgIUIsGhRwp0c2nvU4YSycafPTjzbNcwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwSQXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcNMjEwMzE3MjAzNzEwWhcNMzYwMzE5MDAwMDAwWjB1MUQwQgYDVQQDDDtBcHBsZSBXb3JsZHdpZGUgRGV2ZWxvcGVyIFJlbGF0aW9ucyBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTELMAkGA1UECwwCRzYxEzARBgNVBAoMCkFwcGxlIEluYy4xCzAJBgNVBAYTAlVTMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEbsQKC94PrlWmZXnXgtxzdVJL8T0SGYngDRGpngn3N6PT8JMEb7FDi4bBmPhCnZ3/sq6PF/cGcKXWsL5vOteRhyJ45x3ASP7cOB+aao90fcpxSv/EZFbniAbNgZGhIhpIo4H6MIH3MBIGA1UdEwEB/wQIMAYBAf8CAQAwHwYDVR0jBBgwFoAUu7DeoVgziJqkipnevr3rr9rLJKswRgYIKwYBBQUHAQEEOjA4MDYGCCsGAQUFBzABhipodHRwOi8vb2NzcC5hcHBsZS5jb20vb2NzcDAzLWFwcGxlcm9vdGNhZzMwNwYDVR0fBDAwLjAsoCqgKIYmaHR0cDovL2NybC5hcHBsZS5jb20vYXBwbGVyb290Y2FnMy5jcmwwHQYDVR0OBBYEFD8vlCNR01DJmig97bB85c+lkGKZMA4GA1UdDwEB/wQEAwIBBjAQBgoqhkiG92NkBgIBBAIFADAKBggqhkjOPQQDAwNoADBlAjBAXhSq5IyKogMCPtw490BaB677CaEGJXufQB/EqZGd6CSjiCtOnuMTbXVXmxxcxfkCMQDTSPxarZXvNrkxU3TkUMI33yzvFVVRT4wxWJC994OsdcZ4+RGNsYDyR5gmdr0nDGg="
	REAL_APPLE_SIGNING_CERTIFICATE_BASE64_ENCODED            = "MIIEMTCCA7agAwIBAgIQR8KHzdn554Z/UoradNx9tzAKBggqhkjOPQQDAzB1MUQwQgYDVQQDDDtBcHBsZSBXb3JsZHdpZGUgRGV2ZWxvcGVyIFJlbGF0aW9ucyBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTELMAkGA1UECwwCRzYxEzARBgNVBAoMCkFwcGxlIEluYy4xCzAJBgNVBAYTAlVTMB4XDTI1MDkxOTE5NDQ1MVoXDTI3MTAxMzE3NDcyM1owgZIxQDA+BgNVBAMMN1Byb2QgRUNDIE1hYyBBcHAgU3RvcmUgYW5kIGlUdW5lcyBTdG9yZSBSZWNlaXB0IFNpZ25pbmcxLDAqBgNVBAsMI0FwcGxlIFdvcmxkd2lkZSBEZXZlbG9wZXIgUmVsYXRpb25zMRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNnVvhcv7iT+7Ex5tBMBgrQspHzIsXRi0Yxfek7lv8wEmj/bHiWtNwJqc2BoHzsQiEjP7KFIIKg4Y8y0/nynuAmjggIIMIICBDAMBgNVHRMBAf8EAjAAMB8GA1UdIwQYMBaAFD8vlCNR01DJmig97bB85c+lkGKZMHAGCCsGAQUFBwEBBGQwYjAtBggrBgEFBQcwAoYhaHR0cDovL2NlcnRzLmFwcGxlLmNvbS93d2RyZzYuZGVyMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcC5hcHBsZS5jb20vb2NzcDAzLXd3ZHJnNjAyMIIBHgYDVR0gBIIBFTCCAREwggENBgoqhkiG92NkBQYBMIH+MIHDBggrBgEFBQcCAjCBtgyBs1JlbGlhbmNlIG9uIHRoaXMgY2VydGlmaWNhdGUgYnkgYW55IHBhcnR5IGFzc3VtZXMgYWNjZXB0YW5jZSBvZiB0aGUgdGhlbiBhcHBsaWNhYmxlIHN0YW5kYXJkIHRlcm1zIGFuZCBjb25kaXRpb25zIG9mIHVzZSwgY2VydGlmaWNhdGUgcG9saWN5IGFuZCBjZXJ0aWZpY2F0aW9uIHByYWN0aWNlIHN0YXRlbWVudHMuMDYGCCsGAQUFBwIBFipodHRwOi8vd3d3LmFwcGxlLmNvbS9jZXJ0aWZpY2F0ZWF1dGhvcml0eS8wHQYDVR0OBBYEFIFioG4wMMVA1ku9zJmGNPAVn3eqMA4GA1UdDwEB/wQEAwIHgDAQBgoqhkiG92NkBgsBBAIFADAKBggqhkjOPQQDAwNpADBmAjEA+qXnREC7hXIWVLsLxznjRpIzPf7VHz9V/CTm8+LJlrQepnmcPvGLNcX6XPnlcgLAAjEA5IjNZKgg5pQ79knF4IbTXdKv8vutIDMXDmjPVT3dGvFtsGRwXOywR2kZCdSrfeot"
	LEAF_CERT_PUBLIC_KEY                                     = "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE4rWBxGmFbnPIPQI0zsBKzLxsj8pD\n2vqbr0yPISUx2WQyxmrNql9fhK8YEEyYFV7++p5i4YUSR/o9uQIgCPIhrA==\n-----END PUBLIC KEY-----\n"

	EFFECTIVE_DATE = 1761962975
)

func parsePublicKey(pemStr string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}
	return ecPub, nil
}

func TestValidChainWithoutOCSP(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		LEAF_CERT_BASE64_ENCODED,
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	pubKey, err := cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedPubKey, err := parsePublicKey(LEAF_CERT_PUBLIC_KEY)
	if err != nil {
		t.Fatalf("Failed to parse expected public key: %v", err)
	}
	if pubKey.X.Cmp(expectedPubKey.X) != 0 || pubKey.Y.Cmp(expectedPubKey.Y) != 0 {
		t.Error("Returned public key does not match LEAF_CERT_PUBLIC_KEY")
	}
}

func TestValidChainInvalidIntermediateOIDWithoutOCSP(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		LEAF_CERT_FOR_INTERMEDIATE_CA_INVALID_OID_BASE64_ENCODED,
		INTERMEDIATE_CA_INVALID_OID_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for invalid intermediate OID, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || (vErr.Status != VERIFICATION_FAILURE && vErr.Status != INVALID_CERTIFICATE) {
		t.Errorf("Expected VERIFICATION_FAILURE or INVALID_CERTIFICATE, got %v", err)
	}
}

func TestValidChainInvalidLeafOIDWithoutOCSP(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		LEAF_CERT_INVALID_OID_BASE64_ENCODED,
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for invalid leaf OID, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != VERIFICATION_FAILURE {
		t.Errorf("Expected VERIFICATION_FAILURE, got %v", err)
	}
}

func TestInvalidChainLength(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for invalid chain length, but got none")
	}
	if !strings.Contains(err.Error(), "INVALID_CHAIN_LENGTH") {
		t.Errorf("Expected error message to contain INVALID_CHAIN_LENGTH, got %v", err)
	}
}

func TestInvalidBase64InCertificateList(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		"abc",
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for invalid base64, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != INVALID_CERTIFICATE {
		t.Errorf("Expected INVALID_CERTIFICATE, got %v", err)
	}
}

func TestInvalidDataInCertificateList(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		base64.StdEncoding.EncodeToString([]byte("abc")),
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for invalid data, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != INVALID_CERTIFICATE {
		t.Errorf("Expected INVALID_CERTIFICATE, got %v", err)
	}
}

func TestMalformedRootCert(t *testing.T) {
	malformedRoot := []byte("abc")
	_, err := newChainVerifier([][]byte{malformedRoot})
	if err == nil {
		t.Fatal("Expected error for malformed root during verifier creation, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != INVALID_CERTIFICATE {
		t.Errorf("Expected INVALID_CERTIFICATE, got %v", err)
	}
}

func TestChainDifferentThanRootCertificate(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(REAL_APPLE_ROOT_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		LEAF_CERT_BASE64_ENCODED,
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err == nil {
		t.Fatal("Expected error for mismatching root, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != VERIFICATION_FAILURE {
		t.Errorf("Expected VERIFICATION_FAILURE, got %v", err)
	}
}

func TestValidExpiredChain(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		LEAF_CERT_BASE64_ENCODED,
		INTERMEDIATE_CA_BASE64_ENCODED,
		ROOT_CA_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, false, time.Unix(2280946846, 0))
	if err == nil {
		t.Fatal("Expected error for expired chain, but got none")
	}
	vErr, ok := err.(*VerificationException)
	if !ok || vErr.Status != VERIFICATION_FAILURE {
		t.Errorf("Expected VERIFICATION_FAILURE, got %v", err)
	}
}

func TestAppleChainIsValidWithOCSP(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	rootBytes, _ := base64.StdEncoding.DecodeString(REAL_APPLE_ROOT_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{
		REAL_APPLE_SIGNING_CERTIFICATE_BASE64_ENCODED,
		REAL_APPLE_INTERMEDIATE_BASE64_ENCODED,
		REAL_APPLE_ROOT_BASE64_ENCODED,
	}

	_, err = cv.verifyChain(certs, true, time.Unix(int64(EFFECTIVE_DATE), 0))
	if err != nil {
		t.Fatalf("Expected no error for Apple chain with OCSP, got %v", err)
	}
}

func TestOCSPResponseCaching(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{"cert1", "cert2", "cert3"}
	cacheKey := "cert1|cert2|cert3"

	// Initial hit
	cv.cacheMutex.Lock()
	cv.cache[cacheKey] = cacheEntry{
		publicKey: nil,
		expiry:    time.Now().Add(1 * time.Hour),
	}
	cv.cacheMutex.Unlock()

	pubKey, err := cv.verifyChain(certs, true, time.Now())
	if err != nil {
		t.Fatalf("Expected no error from cache hit, got %v", err)
	}
	if pubKey != nil {
		t.Errorf("Expected nil public key from our mock cache entry, got %v", pubKey)
	}
}

func TestOCSPResponseCachingHasExpiration(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	certs := []string{"cert1", "cert2", "cert3"}
	cacheKey := "cert1|cert2|cert3"

	// Mock entry ready to expire
	cv.cacheMutex.Lock()
	cv.cache[cacheKey] = cacheEntry{
		publicKey: nil,
		expiry:    time.Now().Add(-1 * time.Hour), // Expired
	}
	cv.cacheMutex.Unlock()

	// Should miss cache and fail decoding
	_, err = cv.verifyChain(certs, true, time.Now())
	if err == nil {
		t.Fatal("Expected error for dummy certificates after cache expiration, but got none")
	}
}

func TestOCSPCachingWithDifferentChain(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	chain1 := []string{"leaf1", "int1", "root1"}
	chain2 := []string{"leaf2", "int2", "root2"}

	cv.cacheMutex.Lock()
	cv.cache[strings.Join(chain1, "|")] = cacheEntry{
		publicKey: nil,
		expiry:    time.Now().Add(1 * time.Hour),
	}
	cv.cacheMutex.Unlock()

	// chain1 should hit cache
	_, err = cv.verifyChain(chain1, true, time.Now())
	if err != nil {
		t.Fatalf("Expected no error for chain1 (cache hit), got %v", err)
	}

	// chain2 should NOT hit cache and fail decoding
	_, err = cv.verifyChain(chain2, true, time.Now())
	if err == nil {
		t.Fatal("Expected error for chain2 (cache miss), but got none")
	}
}

func TestOCSPCachingWithSlightlyDifferentChain(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	chain1 := []string{"leaf1", "int1", "root1"}
	chain2 := []string{"leaf1", "int1", "root2"} // Different root

	cv.cacheMutex.Lock()
	cv.cache[strings.Join(chain1, "|")] = cacheEntry{
		publicKey: nil,
		expiry:    time.Now().Add(1 * time.Hour),
	}
	cv.cacheMutex.Unlock()

	// chain1 should hit cache
	_, err = cv.verifyChain(chain1, true, time.Now())
	if err != nil {
		t.Fatalf("Expected no error for chain1 (cache hit), got %v", err)
	}

	// chain2 should NOT hit cache
	_, err = cv.verifyChain(chain2, true, time.Now())
	if err == nil {
		t.Fatal("Expected error for chain2 (cache miss), but got none")
	}
}

func TestCacheEviction(t *testing.T) {
	rootBytes, _ := base64.StdEncoding.DecodeString(ROOT_CA_BASE64_ENCODED)
	cv, err := newChainVerifier([][]byte{rootBytes})
	if err != nil {
		t.Fatalf("Failed to create chain verifier: %v", err)
	}

	// 1. Fill cache to max capacity
	for i := range maxCacheSize {
		key := fmt.Sprintf("chain_%d", i)
		// We can directly manipulate the map for setup because saveToCache enforces limit
		cv.cacheMutex.Lock()
		cv.cache[key] = cacheEntry{
			publicKey: nil, // Value doesn't matter for this test
			expiry:    time.Now().Add(1 * time.Hour),
		}
		cv.cacheMutex.Unlock()
	}

	if len(cv.cache) != maxCacheSize {
		t.Fatalf("Setup failed: cache size expected %d, got %d", maxCacheSize, len(cv.cache))
	}

	// 2. Add one more item - should trigger eviction of a RANDOM item since none are expired
	newItemKey := "new_item_1"
	cv.saveToCache(newItemKey, nil)

	cv.cacheMutex.RLock()
	if len(cv.cache) != maxCacheSize {
		t.Errorf("Eviction failed: cache size expected %d, got %d", maxCacheSize, len(cv.cache))
	}
	if _, ok := cv.cache[newItemKey]; !ok {
		t.Errorf("New item was not added to cache")
	}
	cv.cacheMutex.RUnlock()

	// 3. Test eviction of EXPIRED items
	// First, clear cache and refill
	cv.cacheMutex.Lock()
	// Clear all for fresh start
	for k := range cv.cache {
		delete(cv.cache, k)
	}

	// Fill with max items again
	for i := range maxCacheSize {
		key := fmt.Sprintf("chain_retry_%d", i)
		expiry := time.Now().Add(1 * time.Hour)
		// Mark half as expired
		if i < maxCacheSize/2 {
			expiry = time.Now().Add(-1 * time.Hour)
		}
		cv.cache[key] = cacheEntry{
			publicKey: nil,
			expiry:    expiry,
		}
	}
	cv.cacheMutex.Unlock()

	// Add new item
	newItemKey2 := "new_item_2"
	cv.saveToCache(newItemKey2, nil)

	cv.cacheMutex.RLock()
	// Verify cache size: half were expired/removed (no forced eviction needed).
	// We added one new item, so expected size = (max / 2) + 1.
	expectedSize := (maxCacheSize - (maxCacheSize / 2)) + 1
	if len(cv.cache) != expectedSize {
		t.Errorf("Expired eviction failed: expected size roughly %d, got %d", expectedSize, len(cv.cache))
	}
	if _, ok := cv.cache[newItemKey2]; !ok {
		t.Errorf("New item 2 was not added to cache")
	}

	// Verify no expired items remain
	for k, v := range cv.cache {
		if time.Now().After(v.expiry) {
			t.Errorf("Found expired item in cache: %s", k)
		}
	}
	cv.cacheMutex.RUnlock()
}
