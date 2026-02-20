package appstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	assertEqual(m.t, m.expectedMethod, req.Method, "HTTP method")

	// Validate URL
	expectedParsed, _ := url.Parse(m.expectedURL)
	actualParsed := req.URL

	expectedBase := fmt.Sprintf("%s://%s%s", expectedParsed.Scheme, expectedParsed.Host, expectedParsed.Path)
	actualBase := fmt.Sprintf("%s://%s%s", actualParsed.Scheme, actualParsed.Host, actualParsed.Path)

	assertEqual(m.t, expectedBase, actualBase, "URL base")

	// Validate query parameters
	if len(m.expectedParams) > 0 {
		actualParams := req.URL.Query()
		for key, expectedVals := range m.expectedParams {
			actualVals := actualParams[key]
			assertEqual(m.t, len(expectedVals), len(actualVals), "Param "+key+" values count")
			for i, expectedVal := range expectedVals {
				if i < len(actualVals) {
					assertEqual(m.t, expectedVal, actualVals[i], "Param "+key+" value")
				}
			}
		}
	}

	// Validate headers
	assertTrue(m.t, strings.HasPrefix(req.Header.Get("User-Agent"), "app-store-server-library/go"), "User-Agent header")
	assertEqual(m.t, "application/json", req.Header.Get("Accept"), "Accept header")

	// Validate and decode JWT token
	authHeader := req.Header.Get("Authorization")
	assertTrue(m.t, strings.HasPrefix(authHeader, "Bearer "), "Authorization header")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	header, payload, err := decodeJWTWithoutVerification(tokenString)
	assertNoError(m.t, err, "Failed to decode JWT")

	// Validate JWT claims
	assertEqual(m.t, "appstoreconnect-v1", payload["aud"], "JWT aud")
	assertEqual(m.t, "issuerId", payload["iss"], "JWT iss")
	assertEqual(m.t, "com.example", payload["bid"], "JWT bid")
	assertEqual(m.t, "keyId", header["kid"], "JWT kid")

	// Validate request body
	if m.expectedBinaryData != nil {
		// Binary data validation
		bodyBytes, _ := io.ReadAll(req.Body)
		assertTrue(m.t, bytes.Equal(bodyBytes, m.expectedBinaryData), "Binary body mismatch")
		assertEqual(m.t, m.expectedContentType, req.Header.Get("Content-Type"), "Content-Type header")
	} else if m.expectedBody != nil {
		// JSON body validation - only check that expected fields are present
		var actualBody map[string]any
		bodyBytes, _ := io.ReadAll(req.Body)
		err := json.Unmarshal(bodyBytes, &actualBody)
		assertNoError(m.t, err, "Failed to parse request body")

		// Convert expected body to map
		expectedMap, ok := m.expectedBody.(map[string]any)
		if ok {
			// Check that all expected fields exist with correct values
			for key, expectedVal := range expectedMap {
				actualVal, exists := actualBody[key]
				assertTrue(m.t, exists, "Missing expected field "+key+" in body")
				if exists {
					// Compare values
					expJSON, _ := json.Marshal(expectedVal)
					actJSON, _ := json.Marshal(actualVal)
					assertTrue(m.t, bytes.Equal(expJSON, actJSON), "Field "+key+" mismatch")
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

	// Read response from testdata
	var responseBody []byte
	var err error

	if responseFile != "" {
		responseBody, err = readTestData(fmt.Sprintf("models/%s", responseFile))
		assertNoError(t, err, "Failed to read response file")
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
	assertNoError(t, err, "Failed to read signing key")

	client, err := NewAPIClientWithHTTPClient(signingKey, "keyId", "issuerId", "com.example",
		ENVIRONMENT_LOCAL_TESTING, mockHTTP)
	assertNoError(t, err, "Failed to create API client")

	return client
}

// Test ExtendRenewalDateForAllActiveSubscribers
func TestExtendRenewalDateForAllActiveSubscribers(t *testing.T) {
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
	assertNoError(t, err, "ExtendRenewalDateForAllActiveSubscribers failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "758883e8-151b-47b7-abd0-60c4d804c2f5", response.RequestIdentifier, "requestIdentifier")
}

// Test ExtendSubscriptionRenewalDate
func TestExtendSubscriptionRenewalDate(t *testing.T) {
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
	assertNoError(t, err, "ExtendSubscriptionRenewalDate failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "2312412", response.OriginalTransactionId, "OriginalTransactionId")
	assertEqual(t, "9993", response.WebOrderLineItemId, "WebOrderLineItemId")
	assertEqual(t, true, response.Success, "Success")
	assertEqual(t, Timestamp(1698148900000), response.EffectiveDate, "EffectiveDate")
}

// Test GetAllSubscriptionStatuses
func TestGetAllSubscriptionStatuses(t *testing.T) {
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
	assertNoError(t, err, "GetAllSubscriptionStatuses failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
	assertEqual(t, "com.example", response.BundleId, "BundleId")
	assertEqual(t, int64(5454545), response.AppAppleId, "AppAppleId")

	// Verify subscription group data
	assertEqual(t, 2, len(response.Data), "Subscription groups length")
}

// Test GetRefundHistory
func TestGetRefundHistory(t *testing.T) {
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
	assertNoError(t, err, "GetRefundHistory failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, 2, len(response.SignedTransactions), "SignedTransactions length")
	assertEqual(t, "revision_output", response.Revision, "Revision")
	assertEqual(t, true, response.HasMore, "HasMore")
}

// Test GetStatusOfSubscriptionRenewalDateExtensions
func TestGetStatusOfSubscriptionRenewalDateExtensions(t *testing.T) {
	client := createMockAPIClient(t,
		"getStatusOfSubscriptionRenewalDateExtensionsResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/subscriptions/extend/mass/20fba8a0-2b80-4a7d-a17f-85c1854727f8/com.example.product",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetStatusOfSubscriptionRenewalDateExtensions("20fba8a0-2b80-4a7d-a17f-85c1854727f8", "com.example.product")
	assertNoError(t, err, "GetStatusOfSubscriptionRenewalDateExtensions failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "20fba8a0-2b80-4a7d-a17f-85c1854727f8", response.RequestIdentifier, "RequestIdentifier")
	assertEqual(t, true, response.Complete, "Complete")
	assertEqual(t, Timestamp(1698148900000), response.CompleteDate, "CompleteDate")
	assertEqual(t, int64(30), response.SucceededCount, "SucceededCount")
	assertEqual(t, int64(2), response.FailedCount, "FailedCount")
}

// Test GetTestNotificationStatus
func TestGetTestNotificationStatus(t *testing.T) {
	client := createMockAPIClient(t,
		"getTestNotificationStatusResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/notifications/test/8cd2974c-f905-492a-bf9a-b2f47c791d19",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetTestNotificationStatus("8cd2974c-f905-492a-bf9a-b2f47c791d19")
	assertNoError(t, err, "GetTestNotificationStatus failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "signed_payload", response.SignedPayload, "SignedPayload")
	assertEqual(t, 2, len(response.SendAttempts), "SendAttempts length")
}

// Test GetNotificationHistory
func TestGetNotificationHistory(t *testing.T) {
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
	assertNoError(t, err, "GetNotificationHistory failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "57715481-805a-4283-8499-1c19b5d6b20a", response.PaginationToken, "PaginationToken")
	assertEqual(t, true, response.HasMore, "HasMore")
	assertEqual(t, 2, len(response.NotificationHistory), "NotificationHistory length")
}

// Test GetTransactionHistoryV1
func TestGetTransactionHistoryV1(t *testing.T) {
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
	assertNoError(t, err, "GetTransactionHistory failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "revision_output", response.Revision, "Revision")
	assertEqual(t, true, response.HasMore, "HasMore")
	assertEqual(t, "com.example", response.BundleId, "BundleId")
	assertEqual(t, int64(323232), response.AppAppleId, "AppAppleId")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
}

// Test GetTransactionHistoryV2
func TestGetTransactionHistoryV2(t *testing.T) {
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
	assertNoError(t, err, "GetTransactionHistory failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "revision_output", response.Revision, "Revision")
	assertEqual(t, true, response.HasMore, "HasMore")
	assertEqual(t, "com.example", response.BundleId, "BundleId")
	assertEqual(t, int64(323232), response.AppAppleId, "AppAppleId")
	assertEqual(t, ENVIRONMENT_LOCAL_TESTING, response.Environment, "Environment")
}

// Test GetTransactionInfo
func TestGetTransactionInfo(t *testing.T) {
	client := createMockAPIClient(t,
		"transactionInfoResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/1234",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.GetTransactionInfo("1234")
	assertNoError(t, err, "GetTransactionInfo failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "signed_transaction_info_value", response.SignedTransactionInfo, "SignedTransactionInfo")
}

// Test GetTransactionHistory with unknown environment
func TestGetTransactionHistoryWithUnknownEnvironment(t *testing.T) {
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
	assertNoError(t, err, "GetTransactionHistory failed")

	assertNotNil(t, response, "Response")

	assertEqual(t, Environment("LocalTestingxxx"), response.Environment, "Environment")
	assertEqual(t, false, response.Environment.IsValid(), "Environment.IsValid")
}

// Test LookUpOrderID
func TestLookUpOrderID(t *testing.T) {
	client := createMockAPIClient(t,
		"lookupOrderIdResponse.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/lookup/M12345",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.LookUpOrderID("M12345")
	assertNoError(t, err, "LookUpOrderID failed")

	assertNotNil(t, response, "Response")
	if response.Status != 0 {
		t.Logf("Status: got %v", response.Status)
	}
	assertEqual(t, 2, len(response.SignedTransactions), "SignedTransactions length")
}

// Test RequestTestNotification
func TestRequestTestNotification(t *testing.T) {
	client := createMockAPIClient(t,
		"requestTestNotificationResponse.json",
		"POST",
		"https://local-testing-base-url/inApps/v1/notifications/test",
		map[string][]string{},
		nil,
		200,
	)

	response, err := client.RequestTestNotification()
	assertNoError(t, err, "RequestTestNotification failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "ce3af791-365e-4c60-841b-1674b43c1609", response.TestNotificationToken, "TestNotificationToken")
}

// Test SendConsumptionInformation
func TestSendConsumptionInformation(t *testing.T) {
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
	assertNoError(t, err, "SendConsumptionInformation failed")
}

// Test SetAppAccountToken
func TestSetAppAccountToken(t *testing.T) {
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
	assertNoError(t, err, "SetAppAccountToken failed")
}

// Test error: HTTP 500
func TestAPIError_500(t *testing.T) {
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
	assertError(t, err, "Expected error for 500 status")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 500, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: HTTP 429 Rate Limit
func TestAPIError_429_RateLimit(t *testing.T) {
	client := createMockAPIClient(t,
		"apiTooManyRequestsException.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/1234",
		map[string][]string{},
		nil,
		429,
	)

	_, err := client.GetTransactionInfo("1234")
	assertError(t, err, "Expected error for 429 status")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 429, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assertEqual(t, API_ERROR_RATE_LIMIT_EXCEEDED, *apiErr.APIError, "APIError")
}

// Test error: Invalid Transaction ID
func TestAPIError_InvalidTransactionId(t *testing.T) {
	client := createMockAPIClient(t,
		"invalidTransactionIdError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/invalid",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetTransactionInfo("invalid")
	assertError(t, err, "Expected error for invalid transaction ID")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: Family Transaction Not Supported
func TestAPIError_FamilyTransactionNotSupported(t *testing.T) {
	client := createMockAPIClient(t,
		"familyTransactionNotSupportedError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/subscriptions/9987",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAllSubscriptionStatuses("9987", []Status{})
	assertError(t, err, "Expected error for family transaction")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test error: Invalid App Account Token UUID
func TestAPIError_InvalidAppAccountTokenUUID(t *testing.T) {
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
	assertError(t, err, "Expected error for invalid UUID")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 400, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test UploadImage
func TestUploadImage(t *testing.T) {
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).expectedBinaryData = []byte("fake-image-data")
	client.httpClient.(*MockHTTPClient).expectedContentType = "image/png"
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.UploadImage("img_123", []byte("fake-image-data"))
	assertNoError(t, err, "UploadImage failed")
}

// Test DeleteImage
func TestDeleteImage(t *testing.T) {
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteImage("img_123")
	assertNoError(t, err, "DeleteImage failed")
}

// Test GetImageList
func TestGetImageList(t *testing.T) {
	client := createMockAPIClient(t, "getImageListResponse.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/image/list", nil, nil, 200)

	response, err := client.GetImageList()
	assertNoError(t, err, "GetImageList failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, 1, len(response.ImageIdentifiers), "ImageIdentifiers length")
}

// Test UploadMessage
func TestUploadMessage(t *testing.T) {
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
	assertNoError(t, err, "UploadMessage failed")
}

// Test DeleteMessage
func TestDeleteMessage(t *testing.T) {
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/message/msg_123", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteMessage("msg_123")
	assertNoError(t, err, "DeleteMessage failed")
}

// Test GetMessageList
func TestGetMessageList(t *testing.T) {
	client := createMockAPIClient(t, "getMessageListResponse.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/message/list", nil, nil, 200)

	response, err := client.GetMessageList()
	assertNoError(t, err, "GetMessageList failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, 1, len(response.MessageIdentifiers), "MessageIdentifiers length")
}

// Test ConfigureDefaultMessage
func TestConfigureDefaultMessage(t *testing.T) {
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/inApps/v1/messaging/default/product_1/en-US", nil, map[string]any{
		"messageIdentifier": "msg_123",
	}, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	request := DefaultConfigurationRequest{
		MessageIdentifier: "msg_123",
	}

	err := client.ConfigureDefaultMessage("product_1", "en-US", request)
	assertNoError(t, err, "ConfigureDefaultMessage failed")
}

// Test DeleteDefaultMessage
func TestDeleteDefaultMessage(t *testing.T) {
	client := createMockAPIClient(t, "", "DELETE", "https://local-testing-base-url/inApps/v1/messaging/default/product_1/en-US", nil, nil, 204)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("")

	err := client.DeleteDefaultMessage("product_1", "en-US")
	assertNoError(t, err, "DeleteDefaultMessage failed")
}

// Test GetAppTransactionInfo
func TestGetAppTransactionInfo(t *testing.T) {
	client := createMockAPIClient(t, "appTransactionInfoResponse.json", "GET", "https://local-testing-base-url/inApps/v1/transactions/appTransactions/tx_123", nil, nil, 200)

	response, err := client.GetAppTransactionInfo("tx_123")
	assertNoError(t, err, "GetAppTransactionInfo failed")

	assertNotNil(t, response, "Response")
	assertEqual(t, "signed_app_transaction_info_value", response.SignedAppTransactionInfo, "SignedAppTransactionInfo")
}

// Test APIException.Error
func TestAPIException_Error(t *testing.T) {
	status := API_ERROR_GENERAL_BAD_REQUEST
	e1 := &APIException{
		HTTPStatusCode: 400,
		APIError:       &status,
		ErrorMessage:   "Bad Request",
	}
	expected1 := "API error: 4000000 (code: 400, message: Bad Request)"
	assertEqual(t, expected1, e1.Error(), "e1.Error()")

	e2 := &APIException{
		HTTPStatusCode: 400,
		ErrorMessage:   "Unknown error",
	}
	expected2 := "HTTP error: 400 (message: Unknown error)"
	assertEqual(t, expected2, e2.Error(), "e2.Error()")

	e3 := &APIException{
		HTTPStatusCode: 500,
		ErrorMessage:   "Internal Server Error",
	}
	expected3 := "HTTP error: 500 (message: Internal Server Error)"
	assertEqual(t, expected3, e3.Error(), "e3.Error()")
}

// Test NewAPIClient constructor
func TestNewAPIClient(t *testing.T) {
	signingKey, err := readTestData("certs/testSigningKey.p8")
	assertNoError(t, err, "Failed to read signing key")

	client, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assertNoError(t, err, "Failed to create API client")

	assertNotNil(t, client, "Client")
	assertEqual(t, "https://local-testing-base-url", client.baseURL, "baseURL")
}

// Test GetAppTransactionInfo: Invalid Transaction ID
func TestGetAppTransactionInfo_InvalidTransactionId(t *testing.T) {
	client := createMockAPIClient(t,
		"invalidTransactionIdError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/invalid",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAppTransactionInfo("invalid")
	assertError(t, err, "Expected error for invalid transaction ID")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 400, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assertEqual(t, API_ERROR_INVALID_TRANSACTION_ID, *apiErr.APIError, "APIError")
}

// Test GetAppTransactionInfo: App Transaction Does Not Exist
func TestGetAppTransactionInfo_AppTransactionDoesNotExist(t *testing.T) {
	client := createMockAPIClient(t,
		"appTransactionDoesNotExistError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/nonexistent",
		map[string][]string{},
		nil,
		404,
	)

	_, err := client.GetAppTransactionInfo("nonexistent")
	assertError(t, err, "Expected error for nonexistent app transaction")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 404, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assertEqual(t, API_ERROR_APP_TRANSACTION_DOES_NOT_EXIST_ERROR, *apiErr.APIError, "APIError")
}

// Test GetAppTransactionInfo: Transaction ID Not Found
func TestGetAppTransactionInfo_TransactionIdNotFound(t *testing.T) {
	client := createMockAPIClient(t,
		"transactionIdNotFoundError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/notfound",
		map[string][]string{},
		nil,
		404,
	)

	_, err := client.GetAppTransactionInfo("notfound")
	assertError(t, err, "Expected error for transaction ID not found")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 404, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assertEqual(t, API_ERROR_TRANSACTION_ID_NOT_FOUND, *apiErr.APIError, "APIError")
}

// Test error: Unknown API Error
func TestAPIError_UnknownError(t *testing.T) {
	client := createMockAPIClient(t,
		"apiUnknownError.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/transactions/appTransactions/1234",
		map[string][]string{},
		nil,
		400,
	)

	_, err := client.GetAppTransactionInfo("1234")
	assertError(t, err, "Expected error for unknown API error")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 400, apiErr.HTTPStatusCode, "HTTPStatusCode")
	assertNotNil(t, apiErr.APIError, "APIError for unknown code")
	assertEqual(t, APIError(9990000), *apiErr.APIError, "APIErrorValue")
}

// Test GetTransactionHistory: Malformed App Apple ID
func TestGetTransactionHistory_MalformedAppAppleId(t *testing.T) {
	client := createMockAPIClient(t,
		"transactionHistoryResponseWithMalformedAppAppleId.json",
		"GET",
		"https://local-testing-base-url/inApps/v1/history/1234",
		map[string][]string{},
		nil,
		200,
	)

	_, err := client.GetTransactionHistory("1234", url.Values{}, "", GET_TRANSACTION_HISTORY_VERSION_V1)
	assertError(t, err, "Expected error for malformed AppAppleId")
}

// Test NewAPIClient: Invalid PEM
func TestNewAPIClient_InvalidPEM(t *testing.T) {
	_, err := NewAPIClient([]byte("invalid pem"), "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assertError(t, err, "Expected error for invalid PEM")
	assertEqual(t, "failed to parse PEM block from signing key", err.Error(), "Error message")
}

// Test NewAPIClient: RSA Key (Non-ECDSA)
func TestNewAPIClient_RSAKey(t *testing.T) {
	signingKey, _ := readTestData("certs/rsa_key.pem")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assertError(t, err, "Expected error for RSA key")
	assertEqual(t, "key is not an ECDSA private key", err.Error(), "Error message")
}

// Test NewAPIClient: Non-PKCS8 Key
func TestNewAPIClient_NonPKCS8Key(t *testing.T) {
	signingKey, _ := readTestData("certs/ec_key.pem")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	assertError(t, err, "Expected error for Non-PKCS8 key")
	assertTrue(t, strings.Contains(err.Error(), "failed to parse private key"), "Error message contains failure")
}

// Test NewAPIClient: Invalid Environment
func TestNewAPIClient_InvalidEnvironment(t *testing.T) {
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", Environment("Invalid"))
	assertError(t, err, "Expected error for invalid environment")
}

// Test NewAPIClient: Xcode Environment
func TestNewAPIClient_XcodeEnvironment(t *testing.T) {
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	_, err := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_XCODE)
	assertError(t, err, "Expected error for Xcode environment")
}

// Test APIClient: HTTP Client Error
func TestAPIClient_HTTPClientError(t *testing.T) {
	client := createMockAPIClient(t, "", "GET", "https://local-testing-base-url/inApps/v1/transactions/1234", nil, nil, 0)
	client.httpClient.(*MockHTTPClient).err = errors.New("network error")

	_, err := client.GetTransactionInfo("1234")
	assertError(t, err, "Expected error for network error")
	assertEqual(t, "network error", err.Error(), "Error message")
}

// Test APIClient: Invalid JSON Error Response
func TestAPIClient_InvalidJSONErrorResponse(t *testing.T) {
	client := createMockAPIClient(t, "", "GET", "https://local-testing-base-url/inApps/v1/transactions/1234", nil, nil, 400)
	client.httpClient.(*MockHTTPClient).responseBody = []byte("invalid json")

	_, err := client.GetTransactionInfo("1234")
	assertError(t, err, "Expected error for invalid JSON error response")
	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertTrue(t, strings.Contains(apiErr.ErrorMessage, "failed to decode error response"), "ErrorMessage contains decode failure")
}

// Test UploadMessage with image
func TestUploadMessage_WithImage(t *testing.T) {
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
	assertNoError(t, err, "UploadMessage failed")
}

// Test UploadImage: Error path
func TestUploadImage_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "PUT", "https://local-testing-base-url/inApps/v1/messaging/image/img_123", nil, nil, 500)
	client.httpClient.(*MockHTTPClient).expectedBinaryData = []byte("fake-image-data")
	client.httpClient.(*MockHTTPClient).expectedContentType = "image/png"

	err := client.UploadImage("img_123", []byte("fake-image-data"))
	assertError(t, err, "Expected error for UploadImage")

	apiErr, ok := err.(*APIException)
	assertTrue(t, ok, "Expected APIException")
	assertEqual(t, 500, apiErr.HTTPStatusCode, "HTTPStatusCode")
}

// Test makeRequest: Marshal error
func TestMakeRequest_MarshalError(t *testing.T) {
	signingKey, _ := readTestData("certs/testSigningKey.p8")
	client, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_LOCAL_TESTING)
	// Passing a channel to json.Marshal will cause an error
	err := client.makeRequest("POST", "/test", nil, make(chan int), nil)
	assertError(t, err, "Expected error for JSON marshal failure")
}

// Test NewAPIClient: Standard Environments
func TestNewAPIClient_Environments(t *testing.T) {
	signingKey, _ := readTestData("certs/testSigningKey.p8")

	pClient, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_PRODUCTION)
	assertEqual(t, "https://api.storekit.itunes.apple.com", pClient.baseURL, "Production URL")

	sClient, _ := NewAPIClient(signingKey, "keyId", "issuerId", "com.example", ENVIRONMENT_SANDBOX)
	assertEqual(t, "https://api.storekit-sandbox.itunes.apple.com", sClient.baseURL, "Sandbox URL")
}

// Test internal method: makeRequestWithBinaryBody with destination
func TestMakeRequestWithBinaryBody_WithDestination(t *testing.T) {
	client := createMockAPIClient(t, "", "PUT", "https://local-testing-base-url/test", nil, nil, 200)
	client.httpClient.(*MockHTTPClient).responseBody = []byte(`{"revision":"rev1"}`)

	var response struct {
		Revision string `json:"revision"`
	}
	err := client.makeRequestWithBinaryBody("PUT", "/test", nil, []byte("body"), "text/plain", &response)
	if err != nil {
		assertNoError(t, err, "makeRequestWithBinaryBody failed")
	}
	assertEqual(t, "rev1", response.Revision, "Revision")
}

// Test LookUpOrderID: Error path
func TestLookUpOrderID_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/lookup/order123", nil, nil, 404)
	_, err := client.LookUpOrderID("order123")
	assertError(t, err, "Expected error for LookUpOrderID")
}

// Test RequestTestNotification: Error path
func TestRequestTestNotification_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/test", nil, nil, 500)
	_, err := client.RequestTestNotification()
	assertError(t, err, "Expected error for RequestTestNotification")
}

// Test GetImageList: Error path
func TestGetImageList_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/image/list", nil, nil, 500)
	_, err := client.GetImageList()
	assertError(t, err, "Expected error for GetImageList")
}

// Test GetMessageList: Error path
func TestGetMessageList_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/messaging/message/list", nil, nil, 500)
	_, err := client.GetMessageList()
	assertError(t, err, "Expected error for GetMessageList")
}

// Test ExtendRenewalDateForAllActiveSubscribers: Error path
func TestExtendRenewalDateForAllActiveSubscribers_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/subscriptions/extend/mass", nil, nil, 400)
	request := MassExtendRenewalDateRequest{}
	_, err := client.ExtendRenewalDateForAllActiveSubscribers(request)
	assertError(t, err, "Expected error for ExtendRenewalDateForAllActiveSubscribers")
}

// Test ExtendSubscriptionRenewalDate: Error path
func TestExtendSubscriptionRenewalDate_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "PUT", "https://local-testing-base-url/inApps/v1/subscriptions/extend/tx123", nil, nil, 400)
	request := ExtendRenewalDateRequest{}
	_, err := client.ExtendSubscriptionRenewalDate("tx123", request)
	assertError(t, err, "Expected error for ExtendSubscriptionRenewalDate")
}

// Test GetStatusOfSubscriptionRenewalDateExtensions: Error path
func TestGetStatusOfSubscriptionRenewalDateExtensions_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/subscriptions/extend/mass/req123/prod123", nil, nil, 404)
	_, err := client.GetStatusOfSubscriptionRenewalDateExtensions("req123", "prod123")
	assertError(t, err, "Expected error for GetStatusOfSubscriptionRenewalDateExtensions")
}

// Test GetTestNotificationStatus: Error path
func TestGetTestNotificationStatus_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v1/notifications/test/token123", nil, nil, 404)
	_, err := client.GetTestNotificationStatus("token123")
	assertError(t, err, "Expected error for GetTestNotificationStatus")
}

// Test GetTransactionHistory with revision
func TestGetTransactionHistory_WithRevision(t *testing.T) {
	client := createMockAPIClient(t, "transactionHistoryResponse.json", "GET", "https://local-testing-base-url/inApps/v1/history/1234", map[string][]string{"revision": {"rev1"}}, nil, 200)
	_, err := client.GetTransactionHistory("1234", url.Values{}, "rev1", GET_TRANSACTION_HISTORY_VERSION_V1)
	assertNoError(t, err, "GetTransactionHistory failed")
}

// Test GetRefundHistory with revision
func TestGetRefundHistory_WithRevision(t *testing.T) {
	client := createMockAPIClient(t, "getRefundHistoryResponse.json", "GET", "https://local-testing-base-url/inApps/v2/refund/lookup/1234", map[string][]string{"revision": {"rev1"}}, nil, 200)
	_, err := client.GetRefundHistory("1234", "rev1")
	assertNoError(t, err, "GetRefundHistory failed")
}

// Test GetNotificationHistory with pagination token
func TestGetNotificationHistory_WithToken(t *testing.T) {
	client := createMockAPIClient(t, "getNotificationHistoryResponse.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/history", map[string][]string{"paginationToken": {"token1"}}, nil, 200)
	request := NotificationHistoryRequest{}
	_, err := client.GetNotificationHistory("token1", request)
	assertNoError(t, err, "GetNotificationHistory failed")
}

// Test GetRefundHistory: Error path
func TestGetRefundHistory_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "GET", "https://local-testing-base-url/inApps/v2/refund/lookup/1234", nil, nil, 404)
	_, err := client.GetRefundHistory("1234", "")
	assertError(t, err, "Expected error for GetRefundHistory")
}

// Test GetNotificationHistory: Error path
func TestGetNotificationHistory_Error(t *testing.T) {
	client := createMockAPIClient(t, "apiException.json", "POST", "https://local-testing-base-url/inApps/v1/notifications/history", nil, nil, 500)
	request := NotificationHistoryRequest{}
	_, err := client.GetNotificationHistory("", request)
	assertError(t, err, "Expected error for GetNotificationHistory")
}

// Helper functions
func Int32Ptr(v int32) *int32 {
	return &v
}
