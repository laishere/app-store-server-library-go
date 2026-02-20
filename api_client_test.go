package appstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// MockHTTPClient simulates HTTP responses for testing
type MockHTTPClient struct {
	t                   *testing.T
	expectedMethod      string
	expectedURL         string
	expectedParams      map[string][]string
	expectedBody        any
	expectedContentType string
	expectedBinaryData  []byte
	responseBody        []byte
	responseStatusCode  int
	err                 error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	assert := assert.New(m.t)
	assert.Equal(m.expectedMethod, req.Method, "HTTP method")

	// Validate URL
	expectedParsed, _ := url.Parse(m.expectedURL)
	actualParsed := req.URL

	expectedBase := fmt.Sprintf("%s://%s%s", expectedParsed.Scheme, expectedParsed.Host, expectedParsed.Path)
	actualBase := fmt.Sprintf("%s://%s%s", actualParsed.Scheme, actualParsed.Host, actualParsed.Path)

	assert.Equal(expectedBase, actualBase, "URL base")

	// Validate query parameters
	if len(m.expectedParams) > 0 {
		actualParams := req.URL.Query()
		for key, expectedVals := range m.expectedParams {
			actualVals := actualParams[key]
			assert.Equal(len(expectedVals), len(actualVals), "Param "+key+" values count")
			for i, expectedVal := range expectedVals {
				if i < len(actualVals) {
					assert.Equal(expectedVal, actualVals[i], "Param "+key+" value")
				}
			}
		}
	}

	// Validate headers
	assert.True(strings.HasPrefix(req.Header.Get("User-Agent"), "app-store-server-library/go"), "User-Agent header")
	assert.Equal("application/json", req.Header.Get("Accept"), "Accept header")

	// Validate and decode JWT token
	authHeader := req.Header.Get("Authorization")
	assert.True(strings.HasPrefix(authHeader, "Bearer "), "Authorization header")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	header, payload, err := decodeJWTWithoutVerification(tokenString)
	assert.NoError(err, "Failed to decode JWT")

	// Validate JWT claims
	assert.Equal("appstoreconnect-v1", payload["aud"], "JWT aud")
	assert.Equal("issuerId", payload["iss"], "JWT iss")
	assert.Equal("com.example", payload["bid"], "JWT bid")
	assert.Equal("keyId", header["kid"], "JWT kid")

	// Validate request body
	if m.expectedBinaryData != nil {
		// Binary data validation
		bodyBytes, _ := io.ReadAll(req.Body)
		assert.True(bytes.Equal(bodyBytes, m.expectedBinaryData), "Binary body mismatch")
		assert.Equal(m.expectedContentType, req.Header.Get("Content-Type"), "Content-Type header")
	} else if m.expectedBody != nil {
		// JSON body validation - only check that expected fields are present
		var actualBody map[string]any
		bodyBytes, _ := io.ReadAll(req.Body)
		err := json.Unmarshal(bodyBytes, &actualBody)
		assert.NoError(err, "Failed to parse request body")

		// Convert expected body to map
		expectedMap, ok := m.expectedBody.(map[string]any)
		if ok {
			// Check that all expected fields exist with correct values
			for key, expectedVal := range expectedMap {
				actualVal, exists := actualBody[key]
				assert.True(exists, "Missing expected field "+key+" in body")
				if exists {
					// Compare values
					expJSON, _ := json.Marshal(expectedVal)
					actJSON, _ := json.Marshal(actualVal)
					assert.True(bytes.Equal(expJSON, actJSON), "Field "+key+" mismatch")
				}
			}
		}
	}

	// Create response
	resp := &http.Response{
		StatusCode: m.responseStatusCode,
		Body:       io.NopCloser(bytes.NewReader(m.responseBody)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	if m.err != nil {
		return nil, m.err
	}
	return resp, nil
}

// Helper to create API client with mock HTTP client
func createMockAPIClient(t *testing.T, responseFile string, expectedMethod, expectedURL string,
	expectedParams map[string][]string, expectedBody any, statusCode int) *APIClient {
	assert := assert.New(t)

	// Read response from testdata
	var responseBody []byte
	var err error

	if responseFile != "" {
		responseBody, err = readTestData(fmt.Sprintf("models/%s", responseFile))
		assert.NoError(err, "Failed to read response file")
	} else {
		// Empty response for 204 No Content
		responseBody = []byte("{}")
	}

	mockHTTP := &MockHTTPClient{
		t:                  t,
		expectedMethod:     expectedMethod,
		expectedURL:        expectedURL,
		expectedParams:     expectedParams,
		expectedBody:       expectedBody,
		responseBody:       responseBody,
		responseStatusCode: statusCode,
	}

	signingKey, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	client, err := NewAPIClientWithHTTPClient(signingKey, "keyId", "issuerId", "com.example",
		ENVIRONMENT_LOCAL_TESTING, mockHTTP)
	assert.NoError(err, "Failed to create API client")

	return client
}

// Test ExtendRenewalDateForAllActiveSubscribers
func TestExtendRenewalDateForAllActiveSubscribers(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"extendRenewalDateForAllActiveSubscribersResponse.json",
		"POST",
		"https://local-testing-base-url/inApps/v1/subscriptions/extend/mass",
		map[string][]string{},
		map[string]any{
			"extendByDays":           float64(45),
			"extendReasonCode":       float64(1),
			"requestIdentifier":      "fdf964a4-233b-486c-aac1-97d8d52688ac",
			"storefrontCountryCodes": []any{"USA", "MEX"},
			"productId":              "com.example.productId",
		},
		200,
	)

	request := MassExtendRenewalDateRequest{
		ExtendByDays:           45,
		ExtendReasonCode:       EXTEND_REASON_CODE_CUSTOMER_SATISFACTION,
		RequestIdentifier:      "fdf964a4-233b-486c-aac1-97d8d52688ac",
		StorefrontCountryCodes: []string{"USA", "MEX"},
		ProductId:              "com.example.productId",
	}

	response, err := client.ExtendRenewalDateForAllActiveSubscribers(request)
	assert.NoError(err, "ExtendRenewalDateForAllActiveSubscribers failed")

	assert.NotNil(response, "Response")
	assert.Equal("758883e8-151b-47b7-abd0-60c4d804c2f5", response.RequestIdentifier, "requestIdentifier")
}

// Test ExtendSubscriptionRenewalDate
func TestExtendSubscriptionRenewalDate(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"extendSubscriptionRenewalDateResponse.json",
		"PUT",
		"https://local-testing-base-url/inApps/v1/subscriptions/extend/4124214",
		map[string][]string{},
		map[string]any{
			"extendByDays":      float64(45),
			"extendReasonCode":  float64(1),
			"requestIdentifier": "fdf964a4-233b-486c-aac1-97d8d52688ac",
		},
		200,
	)

	request := ExtendRenewalDateRequest{
		ExtendByDays:      45,
		ExtendReasonCode:  EXTEND_REASON_CODE_CUSTOMER_SATISFACTION,
		RequestIdentifier: "fdf964a4-233b-486c-aac1-97d8d52688ac",
	}

	response, err := client.ExtendSubscriptionRenewalDate("4124214", request)
	assert.NoError(err, "ExtendSubscriptionRenewalDate failed")

	assert.NotNil(response, "Response")
	assert.Equal("2312412", response.OriginalTransactionId, "OriginalTransactionId")
	assert.Equal("9993", response.WebOrderLineItemId, "WebOrderLineItemId")
	assert.Equal(true, response.Success, "Success")
	assert.Equal(Timestamp(1698148900000), response.EffectiveDate, "EffectiveDate")
}

// Test GetAllSubscriptionStatuses
func TestGetAllSubscriptionStatuses(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"getAllSubscriptionStatusesResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/subscriptions/4321",
		map[string][]string{
			"status": {"2", "1"},
		},
		nil,
		200,
	)

	response, err := client.GetAllSubscriptionStatuses("4321", []Status{
		STATUS_EXPIRED,
		STATUS_ACTIVE,
	})
	assert.NoError(err, "GetAllSubscriptionStatuses failed")

	assert.NotNil(response, "Response")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
	assert.Equal("com.example", response.BundleId, "BundleId")
	assert.Equal(int64(5454545), response.AppAppleId, "AppAppleId")

	// Verify subscription group data
	assert.Equal(2, len(response.Data), "Subscription groups length")
}

// Test GetRefundHistory
func TestGetRefundHistory(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"getRefundHistoryResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v2/refund/lookup/555555",
		map[string][]string{
			"revision": {"revision_input"},
		},
		nil,
		200,
	)

	response, err := client.GetRefundHistory("555555", "revision_input")
	assert.NoError(err, "GetRefundHistory failed")

	assert.NotNil(response, "Response")
	assert.Equal(2, len(response.SignedTransactions), "SignedTransactions length")
	assert.Equal("revision_output", response.Revision, "Revision")
	assert.Equal(true, response.HasMore, "HasMore")
}

// Test GetStatusOfSubscriptionRenewalDateExtensions
func TestGetStatusOfSubscriptionRenewalDateExtensions(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"getStatusOfSubscriptionRenewalDateExtensionsResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/subscriptions/extend/mass/20fba8a0-2b80-4a7d-a17f-85c1854727f8/com.example.product",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetStatusOfSubscriptionRenewalDateExtensions("20fba8a0-2b80-4a7d-a17f-85c1854727f8", "com.example.product")
	assert.NoError(err, "GetStatusOfSubscriptionRenewalDateExtensions failed")

	assert.NotNil(response, "Response")
	assert.Equal("20fba8a0-2b80-4a7d-a17f-85c1854727f8", response.RequestIdentifier, "RequestIdentifier")
	assert.Equal(true, response.Complete, "Complete")
	assert.Equal(Timestamp(1698148900000), response.CompleteDate, "CompleteDate")
	assert.Equal(int64(30), response.SucceededCount, "SucceededCount")
	assert.Equal(int64(2), response.FailedCount, "FailedCount")
}

// Test GetTestNotificationStatus
func TestGetTestNotificationStatus(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"getTestNotificationStatusResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/notifications/test/8cd2974c-f905-492a-bf9a-b2f47c791d19",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetTestNotificationStatus("8cd2974c-f905-492a-bf9a-b2f47c791d19")
	assert.NoError(err, "GetTestNotificationStatus failed")

	assert.NotNil(response, "Response")
	assert.Equal("signed_payload", response.SignedPayload, "SignedPayload")
	assert.Equal(2, len(response.SendAttempts), "SendAttempts length")
}

// Test GetNotificationHistory
func TestGetNotificationHistory(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"getNotificationHistoryResponse.json",
		"POST",
		"https://local-testing-base-url/inApps/v1/notifications/history",
		map[string][]string{
			"paginationToken": {"a036bc0e-52b8-4bee-82fc-8c24cb6715d6"},
		},
		map[string]any{
			"startDate":           float64(1698148900000),
			"endDate":             float64(1698148950000),
			"notificationType":    "SUBSCRIBED",
			"notificationSubtype": "INITIAL_BUY",
			"transactionId":       "999733843",
			"onlyFailures":        true,
		},
		200,
	)

	request := NotificationHistoryRequest{
		StartDate:           Timestamp(1698148900000),
		EndDate:             Timestamp(1698148950000),
		NotificationType:    NOTIFICATION_TYPE_SUBSCRIBED,
		NotificationSubtype: SUBTYPE_INITIAL_BUY,
		TransactionId:       "999733843",
		OnlyFailures:        true,
	}

	response, err := client.GetNotificationHistory("a036bc0e-52b8-4bee-82fc-8c24cb6715d6", request)
	assert.NoError(err, "GetNotificationHistory failed")

	assert.NotNil(response, "Response")
	assert.Equal("57715481-805a-4283-8499-1c19b5d6b20a", response.PaginationToken, "PaginationToken")
	assert.Equal(true, response.HasMore, "HasMore")
	assert.Equal(2, len(response.NotificationHistory), "NotificationHistory length")
}

// Test GetTransactionHistoryV1
func TestGetTransactionHistoryV1(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionHistoryResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/history/1234",
		map[string][]string{
			"revision":                    {"revision_input"},
			"startDate":                   {"123455"},
			"endDate":                     {"123456"},
			"productId":                   {"com.example.1", "com.example.2"},
			"productType":                 {"CONSUMABLE", "AUTO_RENEWABLE"},
			"sort":                        {"ASCENDING"},
			"subscriptionGroupIdentifier": {"sub_group_id", "sub_group_id_2"},
			"inAppOwnershipType":          {"FAMILY_SHARED"},
			"revoked":                     {"false"},
		},
		nil,
		200,
	)

	queryParams := url.Values{}
	queryParams.Set("revision", "revision_input")
	queryParams.Set("startDate", "123455")
	queryParams.Set("endDate", "123456")
	queryParams.Add("productId", "com.example.1")
	queryParams.Add("productId", "com.example.2")
	queryParams.Add("productType", "CONSUMABLE")
	queryParams.Add("productType", "AUTO_RENEWABLE")
	queryParams.Set("sort", "ASCENDING")
	queryParams.Add("subscriptionGroupIdentifier", "sub_group_id")
	queryParams.Add("subscriptionGroupIdentifier", "sub_group_id_2")
	queryParams.Set("inAppOwnershipType", "FAMILY_SHARED")
	queryParams.Set("revoked", "false")

	response, err := client.GetTransactionHistory("1234", queryParams, "", GET_TRANSACTION_HISTORY_VERSION_V1)
	assert.NoError(err, "GetTransactionHistory failed")

	assert.NotNil(response, "Response")
	assert.Equal("revision_output", response.Revision, "Revision")
	assert.Equal(true, response.HasMore, "HasMore")
	assert.Equal("com.example", response.BundleId, "BundleId")
	assert.Equal(int64(323232), response.AppAppleId, "AppAppleId")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
}

// Test GetTransactionHistoryV2
func TestGetTransactionHistoryV2(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionHistoryResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v2/history/1234",
		map[string][]string{
			"revision":                    {"revision_input"},
			"startDate":                   {"123455"},
			"endDate":                     {"123456"},
			"productId":                   {"com.example.1", "com.example.2"},
			"productType":                 {"CONSUMABLE", "AUTO_RENEWABLE"},
			"sort":                        {"ASCENDING"},
			"subscriptionGroupIdentifier": {"sub_group_id", "sub_group_id_2"},
			"inAppOwnershipType":          {"FAMILY_SHARED"},
			"revoked":                     {"false"},
		},
		nil,
		200,
	)

	queryParams := url.Values{}
	queryParams.Set("revision", "revision_input")
	queryParams.Set("startDate", "123455")
	queryParams.Set("endDate", "123456")
	queryParams.Add("productId", "com.example.1")
	queryParams.Add("productId", "com.example.2")
	queryParams.Add("productType", "CONSUMABLE")
	queryParams.Add("productType", "AUTO_RENEWABLE")
	queryParams.Set("sort", "ASCENDING")
	queryParams.Add("subscriptionGroupIdentifier", "sub_group_id")
	queryParams.Add("subscriptionGroupIdentifier", "sub_group_id_2")
	queryParams.Set("inAppOwnershipType", "FAMILY_SHARED")
	queryParams.Set("revoked", "false")

	response, err := client.GetTransactionHistory("1234", queryParams, "", GET_TRANSACTION_HISTORY_VERSION_V2)
	assert.NoError(err, "GetTransactionHistory failed")

	assert.NotNil(response, "Response")
	assert.Equal("revision_output", response.Revision, "Revision")
	assert.Equal(true, response.HasMore, "HasMore")
	assert.Equal("com.example", response.BundleId, "BundleId")
	assert.Equal(int64(323232), response.AppAppleId, "AppAppleId")
	assert.Equal(ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
}

// Test GetTransactionInfo
func TestGetTransactionInfo(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionInfoResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/1234",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetTransactionInfo("1234")
	assert.NoError(err, "GetTransactionInfo failed")

	assert.NotNil(response, "Response")
	assert.Equal("signed_transaction_info_value", response.SignedTransactionInfo, "SignedTransactionInfo")
}

// Test GetTransactionHistory with unknown environment
func TestGetTransactionHistoryWithUnknownEnvironment(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionHistoryResponseWithMalformedEnvironment.json",
		"GET",
		"https://local-testing-base-url/inApps/v2/history/1234",
		map[string][]string{
			"revision": {"revision_input"},
		},
		nil,
		200,
	)

	response, err := client.GetTransactionHistory("1234", url.Values{"revision": {"revision_input"}}, "", GET_TRANSACTION_HISTORY_VERSION_V2)
	assert.NoError(err, "GetTransactionHistory failed")

	assert.NotNil(response, "Response")

	assert.Equal(Environment("LocalTestingxxx"), response.Environment, "Environment")
	assert.Equal(false, response.Environment.IsValid(), "Environment.IsValid")
}

// Test LookUpOrderID
func TestLookUpOrderID(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"lookupOrderIdResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/lookup/M12345",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.LookUpOrderID("M12345")
	assert.NoError(err, "LookUpOrderID failed")

	assert.NotNil(response, "Response")
	if response.Status != 0 {
		t.Logf("Status: got %v", response.Status)
	}
	assert.Equal(2, len(response.SignedTransactions), "SignedTransactions length")
}

// Test RequestTestNotification
func TestRequestTestNotification(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"requestTestNotificationResponse.json",
		"POST",
		"https://local-testing-base-url/inApps/v1/notifications/test",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.RequestTestNotification()
	assert.NoError(err, "RequestTestNotification failed")

	assert.NotNil(response, "Response")
	assert.Equal("ce3af791-365e-4c60-841b-1674b43c1609", response.TestNotificationToken, "TestNotificationToken")
}

// Test SendConsumptionInformation
func TestSendConsumptionInformation(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"",
		"PUT",
		"https://local-testing-base-url/inApps/v2/transactions/consumption/49571273",
		map[string][]string{},
		map[string]any{
			"customerConsented":     true,
			"sampleContentProvided": false,
			"consumptionPercentage": float64(50000),
			"deliveryStatus":        float64(0),
			"refundPreference":      float64(1),
		},
		204,
	)

	// Manually set empty response for 204 No Content
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := ConsumptionRequest{
		CustomerConsented:     true,
		SampleContentProvided: false,
		DeliveryStatus:        DELIVERY_STATUS_DELIVERED_AND_WORKING_PROPERLY,
		ConsumptionPercentage: 50000,
		RefundPreference:      REFUND_PREFERENCE_PREFER_REFUND,
	}

	err := client.SendConsumptionInformation("49571273", request)
	assert.NoError(err, "SendConsumptionInformation failed")
}

// Test SetAppAccountToken
func TestSetAppAccountToken(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"",
		"PUT",
		"https://local-testing-base-url/inApps/v1/transactions/1234/appAccountToken",
		map[string][]string{},
		map[string]any{
			"appAccountToken": "7e3fb20b-4cdb-47cc-936d-99d65f608138",
		},
		204,
	)

	// Manually set empty response for 204 No Content
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := UpdateAppAccountTokenRequest{
		AppAccountToken: "7e3fb20b-4cdb-47cc-936d-99d65f608138",
	}

	err := client.SetAppAccountToken("1234", request)
	assert.NoError(err, "SetAppAccountToken failed")
}

// Test error: HTTP 500
func TestAPIError_500(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/1234",
		map[string][]string{},
		nil,
		500,
	)

	// Set error response
	errorBody := map[string]any{
		"errorCode":    5000000,
		"errorMessage": "An error occurred",
	}
	errorJSON, _ := json.Marshal(errorBody)
	client.httpClient.(*MockHTTPClient).responseBody = errorJSON

	_, err := client.GetTransactionInfo("1234")
	assert.Error(err, "Expected error for 500 status")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(500, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: HTTP 429 Rate Limit
func TestAPIError_429_RateLimit(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"apiTooManyRequestsException.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/1234",
		map[string][]string{},
		nil,
		429,
	)

	_, err := client.GetTransactionInfo("1234")
	assert.Error(err, "Expected error for 429 status")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(429, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assert.Equal(API_ERROR_RATE_LIMIT_EXCEEDED, *apiErr.APIError, "APIError")
}

// Test error: Invalid Transaction ID
func TestAPIError_InvalidTransactionId(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"invalidTransactionIdError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/invalid",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetTransactionInfo("invalid")
	assert.Error(err, "Expected error for invalid transaction ID")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: Family Transaction Not Supported
func TestAPIError_FamilyTransactionNotSupported(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"familyTransactionNotSupportedError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/subscriptions/9987",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAllSubscriptionStatuses("9987", []Status{})
	assert.Error(err, "Expected error for family transaction")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: Invalid App Account Token UUID
func TestAPIError_InvalidAppAccountTokenUUID(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"invalidAppAccountTokenUUIDError.json",
		"PUT",
		"https://local-testing-base-url/inApps/v2/transactions/consumption/1234",
		map[string][]string{},
		map[string]any{
			"customerConsented":     true,
			"sampleContentProvided": false,
			"deliveryStatus":        float64(0),
		},
		400,
	)

	request := ConsumptionRequest{
		CustomerConsented:     true,
		SampleContentProvided: false,
		DeliveryStatus:        DELIVERY_STATUS_DELIVERED_AND_WORKING_PROPERLY,
	}

	err := client.SendConsumptionInformation("1234", request)
	assert.Error(err, "Expected error for invalid UUID")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test UploadImage
func TestUploadImage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).expectedBinaryData = []byte("fake-image-data")
	client.httpClient.(*MockHTTPClient).expectedContentType = "image/png"
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.UploadImage("img_123", []byte("fake-image-data"))
	assert.NoError(err, "UploadImage failed")
}

// Test DeleteImage
func TestDeleteImage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteImage("img_123")
	assert.NoError(err, "DeleteImage failed")
}

// Test GetImageList
func TestGetImageList(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "getImageListResponse.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/image/list", nil, nil, 200)

	response, err := client.GetImageList()
	assert.NoError(err, "GetImageList failed")

	assert.NotNil(response, "Response")
	assert.Equal(1, len(response.ImageIdentifiers), "ImageIdentifiers length")
}

// Test UploadMessage
func TestUploadMessage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/message/msg_123", nil, map[string]any{
		"header": "Hello",
		"body":   "World",
	}, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := UploadMessageRequestBody{
		Header: "Hello",
		Body:   "World",
	}

	err := client.UploadMessage("msg_123", request)
	assert.NoError(err, "UploadMessage failed")
}

// Test DeleteMessage
func TestDeleteMessage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/message/msg_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteMessage("msg_123")
	assert.NoError(err, "DeleteMessage failed")
}

// Test GetMessageList
func TestGetMessageList(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "getMessageListResponse.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/message/list", nil, nil, 200)

	response, err := client.GetMessageList()
	assert.NoError(err, "GetMessageList failed")

	assert.NotNil(response, "Response")
	assert.Equal(1, len(response.MessageIdentifiers), "MessageIdentifiers length")
}

// Test ConfigureDefaultMessage
func TestConfigureDefaultMessage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/default/product_1/en-US", nil, map[string]any{
		"messageIdentifier": "msg_123",
	}, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := DefaultConfigurationRequest{
		MessageIdentifier: "msg_123",
	}

	err := client.ConfigureDefaultMessage("product_1", "en-US", request)
	assert.NoError(err, "ConfigureDefaultMessage failed")
}

// Test DeleteDefaultMessage
func TestDeleteDefaultMessage(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/default/product_1/en-US", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteDefaultMessage("product_1", "en-US")
	assert.NoError(err, "DeleteDefaultMessage failed")
}

// Test GetAppTransactionInfo
func TestGetAppTransactionInfo(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "appTransactionInfoResponse.json", "GET", "https://local-testing-base-url/inApps/v1/transactions/appTransactions/tx_123", nil, nil, 200)

	response, err := client.GetAppTransactionInfo("tx_123")
	assert.NoError(err, "GetAppTransactionInfo failed")

	assert.NotNil(response, "Response")
	assert.Equal("signed_app_transaction_info_value", response.SignedAppTransactionInfo, "SignedAppTransactionInfo")
}

// Test APIException.Error
func TestAPIException_Error(t *testing.T) {
	assert := assert.New(t)
	status := API_ERROR_GENERAL_BAD_REQUEST
	e1 := &APIException{
		HTTPStatusCode: 400,
		APIError:       &status,
		ErrorMessage:   "Bad Request",
	}
	expected1 := "API error: 4000000 (code: 400, message: Bad Request)"
	assert.Equal(expected1, e1.Error(), "e1.Error()")

	e2 := &APIException{
		HTTPStatusCode: 400,
		ErrorMessage:   "Unknown error",
	}
	expected2 := "HTTP error: 400 (message: Unknown error)"
	assert.Equal(expected2, e2.Error(), "e2.Error()")

	e3 := &APIException{
		HTTPStatusCode: 500,
		ErrorMessage:   "Internal Server Error",
	}
	expected3 := "HTTP error: 500 (message: Internal Server Error)"
	assert.Equal(expected3, e3.Error(), "e3.Error()")
}

// Test NewAPIClient constructor
func TestNewAPIClient(t *testing.T) {
	assert := assert.New(t)
	signingKey, err := readTestData("certs/testSigningKey.p8")
	assert.NoError(err, "Failed to read signing key")

	client, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assert.NoError(err, "Failed to create API client")

	assert.NotNil(client, "Client")
	assert.Equal("https://local-testing-base-url", client.baseURL, "baseURL")
}

// Test GetAppTransactionInfo: Invalid Transaction ID
func TestGetAppTransactionInfo_InvalidTransactionId(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"invalidTransactionIdError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/invalid",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAppTransactionInfo("invalid")
	assert.Error(err, "Expected error for invalid transaction ID")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(400, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assert.Equal(API_ERROR_INVALID_TRANSACTION_ID, *apiErr.APIError, "APIError")
}

// Test GetAppTransactionInfo: App Transaction Does Not Exist
func TestGetAppTransactionInfo_AppTransactionDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"appTransactionDoesNotExistError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/nonexistent",
		map[string][]string{},
		nil,
		404,
	)

	_, err := client.GetAppTransactionInfo("nonexistent")
	assert.Error(err, "Expected error for nonexistent app transaction")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(404, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assert.Equal(API_ERROR_APP_TRANSACTION_DOES_NOT_EXIST_ERROR, *apiErr.APIError, "APIError")
}

// Test GetAppTransactionInfo: Transaction ID Not Found
func TestGetAppTransactionInfo_TransactionIdNotFound(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionIdNotFoundError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/notfound",
		map[string][]string{},
		nil,
		404,
	)

	_, err := client.GetAppTransactionInfo("notfound")
	assert.Error(err, "Expected error for transaction ID not found")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(404, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assert.Equal(API_ERROR_TRANSACTION_ID_NOT_FOUND, *apiErr.APIError, "APIError")
}

// Test error: Unknown API Error
func TestAPIError_UnknownError(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"apiUnknownError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/1234",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAppTransactionInfo("1234")
	assert.Error(err, "Expected error for unknown API error")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(400, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assert.NotNil(apiErr.APIError, "APIError for unknown code")
	assert.Equal(APIError(9990000), *apiErr.APIError, "APIErrorValue")
}

// Test GetTransactionHistory: Malformed App Apple ID
func TestGetTransactionHistory_MalformedAppAppleId(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t,
		"transactionHistoryResponseWithMalformedAppAppleId.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/history/1234",
		map[string][]string{},
		nil,
		200,
	)

	_, err := client.GetTransactionHistory("1234", url.Values{}, "", GET_TRANSACTION_HISTORY_VERSION_V1)
	assert.Error(err, "Expected error for malformed AppAppleId")
}

// Test NewAPIClient: Invalid PEM
func TestNewAPIClient_InvalidPEM(t *testing.T) {
	assert := assert.New(t)
	_, err := NewAPIClient([]byte("invalid pem"), "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assert.Error(err, "Expected error for invalid PEM")
	assert.Equal("failed to parse PEM block from signing key", err.Error(), "Error message")
}

// Test NewAPIClient: RSA Key (Non-ECDSA)
func TestNewAPIClient_RSAKey(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/rsa_key.pem")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assert.Error(err, "Expected error for RSA key")
	assert.Equal("key is not an ECDSA private key", err.Error(), "Error message")
}

// Test NewAPIClient: Non-PKCS8 Key
func TestNewAPIClient_NonPKCS8Key(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/ec_key.pem")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assert.Error(err, "Expected error for Non-PKCS8 key")
	assert.True(strings.Contains(err.Error(), "failed to parse private key"), "Error message contains failure")
}

// Test NewAPIClient: Invalid Environment
func TestNewAPIClient_InvalidEnvironment(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", Environment("Invalid"))
	assert.Error(err, "Expected error for invalid environment")
}

// Test NewAPIClient: Xcode Environment
func TestNewAPIClient_XcodeEnvironment(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_XCODE)
	assert.Error(err, "Expected error for Xcode environment")
}

// Test APIClient: HTTP Client Error
func TestAPIClient_HTTPClientError(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "GET", "https://local-testing-base-url/inApps/v1/transactions/1234", nil, nil, 0)
	client.httpClient.(*MockHTTPClient).err = errors.New("network error")

	_, err := client.GetTransactionInfo("1234")
	assert.Error(err, "Expected error for network error")
	assert.Equal("network error", err.Error(), "Error message")
}

// Test APIClient: Invalid JSON Error Response
func TestAPIClient_InvalidJSONErrorResponse(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "GET", "https://local-testing-base-url/inApps/v1/transactions/1234", nil, nil, 400)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("invalid json")

	_, err := client.GetTransactionInfo("1234")
	assert.Error(err, "Expected error for invalid JSON error response")
	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.True(strings.Contains(apiErr.ErrorMessage, "failed to decode error response"), "ErrorMessage contains decode failure")
}

// Test UploadMessage with image
func TestUploadMessage_WithImage(t *testing.T) {
	assert := assert.New(t)
	expectedBody := map[string]any{
		"header": "Hello",
		"body":   "World",
		"image": map[string]any{
			"imageIdentifier": "img_123",
			"altText":         "Alt text",
		},
	}
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/message/msg_123", nil, expectedBody, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := UploadMessageRequestBody{
		Header: "Hello",
		Body:   "World",
		Image: &UploadMessageImage{
			ImageIdentifier: "img_123",
			AltText:         "Alt text",
		},
	}

	err := client.UploadMessage("msg_123", request)
	assert.NoError(err, "UploadMessage failed")
}

// Test UploadImage: Error path
func TestUploadImage_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "PUT", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 500)
	client.httpClient.(*MockHTTPClient).expectedBinaryData = []byte("fake-image-data")
	client.httpClient.(*MockHTTPClient).expectedContentType = "image/png"

	err := client.UploadImage("img_123", []byte("fake-image-data"))
	assert.Error(err, "Expected error for UploadImage")

	apiErr, ok := err.(*APIException)
	assert.True(ok, "Expected APIException")
	assert.Equal(500, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test makeRequest: Marshal error
func TestMakeRequest_MarshalError(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	client, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	// Passing a channel to json.Marshal will cause an error
	err := client.makeRequest("POST", "/test", nil, make(chan int), nil)
	assert.Error(err, "Expected error for JSON marshal failure")
}

// Test NewAPIClient: Standard Environments
func TestNewAPIClient_Environments(t *testing.T) {
	assert := assert.New(t)
	signingKey, _ := readTestData("certs/testSigningKey.p8")

	pClient, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_PRODUCTION)
	assert.Equal("https://api.storekit.itunes.apple.com", pClient.baseURL, "Production URL")

	sClient, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_SANDBOX)
	assert.Equal("https://api.storekit-sandbox.itunes.apple.com", sClient.baseURL, "Sandbox URL")
}

// Test internal method: makeRequestWithBinaryBody with destination
func TestMakeRequestWithBinaryBody_WithDestination(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/test", nil, nil, 200)
	client.httpClient.(*MockHTTPClient).responseBody = []byte(`{"revision":"rev1"}`)

	var response struct {
		Revision string `json:"revision"`
	}
	err := client.makeRequestWithBinaryBody("PUT", "/test", nil, []byte("body"), "text/plain", &response)
	if err != nil {
		assert.NoError(err, "makeRequestWithBinaryBody failed")
	}
	assert.Equal("rev1", response.Revision, "Revision")
}

// Test LookUpOrderID: Error path
func TestLookUpOrderID_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/lookup/order123", nil, nil, 404)
	_, err := client.LookUpOrderID("order123")
	assert.Error(err, "Expected error for LookUpOrderID")
}

// Test RequestTestNotification: Error path
func TestRequestTestNotification_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/test", nil, nil, 500)
	_, err := client.RequestTestNotification()
	assert.Error(err, "Expected error for RequestTestNotification")
}

// Test GetImageList: Error path
func TestGetImageList_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/image/list", nil, nil, 500)
	_, err := client.GetImageList()
	assert.Error(err, "Expected error for GetImageList")
}

// Test GetMessageList: Error path
func TestGetMessageList_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/message/list", nil, nil, 500)
	_, err := client.GetMessageList()
	assert.Error(err, "Expected error for GetMessageList")
}

// Test ExtendRenewalDateForAllActiveSubscribers: Error path
func TestExtendRenewalDateForAllActiveSubscribers_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/subscriptions/extend/mass", nil, nil, 400)
	request := MassExtendRenewalDateRequest{}
	_, err := client.ExtendRenewalDateForAllActiveSubscribers(request)
	assert.Error(err, "Expected error for ExtendRenewalDateForAllActiveSubscribers")
}

// Test ExtendSubscriptionRenewalDate: Error path
func TestExtendSubscriptionRenewalDate_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "PUT", "https://local-testing-base-url/inApps/v1/subscriptions/extend/tx123", nil, nil, 400)
	request := ExtendRenewalDateRequest{}
	_, err := client.ExtendSubscriptionRenewalDate("tx123", request)
	assert.Error(err, "Expected error for ExtendSubscriptionRenewalDate")
}

// Test GetStatusOfSubscriptionRenewalDateExtensions: Error path
func TestGetStatusOfSubscriptionRenewalDateExtensions_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/subscriptions/extend/mass/req123/prod123", nil, nil, 404)
	_, err := client.GetStatusOfSubscriptionRenewalDateExtensions("req123", "prod123")
	assert.Error(err, "Expected error for GetStatusOfSubscriptionRenewalDateExtensions")
}

// Test GetTestNotificationStatus: Error path
func TestGetTestNotificationStatus_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/notifications/test/token123", nil, nil, 404)
	_, err := client.GetTestNotificationStatus("token123")
	assert.Error(err, "Expected error for GetTestNotificationStatus")
}

// Test GetTransactionHistory with revision
func TestGetTransactionHistory_WithRevision(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "transactionHistoryResponse.json", "GET", "https://local-testing-base-url/inApps/v1/history/1234", map[string][]string{"revision": {"rev1"}}, nil, 200)
	_, err := client.GetTransactionHistory("1234", url.Values{}, "rev1", GET_TRANSACTION_HISTORY_VERSION_V1)
	assert.NoError(err, "GetTransactionHistory failed")
}

// Test GetRefundHistory with revision
func TestGetRefundHistory_WithRevision(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "getRefundHistoryResponse.json", "GET", "https://local-testing-base-url/inApps/v2/refund/lookup/1234", map[string][]string{"revision": {"rev1"}}, nil, 200)
	_, err := client.GetRefundHistory("1234", "rev1")
	assert.NoError(err, "GetRefundHistory failed")
}

// Test GetNotificationHistory with pagination token
func TestGetNotificationHistory_WithToken(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "getNotificationHistoryResponse.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/history", map[string][]string{"paginationToken": {"token1"}}, nil, 200)
	request := NotificationHistoryRequest{}
	_, err := client.GetNotificationHistory("token1", request)
	assert.NoError(err, "GetNotificationHistory failed")
}

// Test GetRefundHistory: Error path
func TestGetRefundHistory_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v2/refund/lookup/1234", nil, nil, 404)
	_, err := client.GetRefundHistory("1234", "")
	assert.Error(err, "Expected error for GetRefundHistory")
}

// Test GetNotificationHistory: Error path
func TestGetNotificationHistory_Error(t *testing.T) {
	assert := assert.New(t)
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/history", nil, nil, 500)
	request := NotificationHistoryRequest{}
	_, err := client.GetNotificationHistory("", request)
	assert.Error(err, "Expected error for GetNotificationHistory")
}

// Helper functions
func Int32Ptr(v int32) *int32 {
	return &v
}
