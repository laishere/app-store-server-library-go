package appstore

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// HTTPClient is an interface for making HTTP requests.
// It allows for custom HTTP client implementations and facilitates testing with mock clients.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIException represents an error response from the App Store Server API.
// It contains the HTTP status code, the parsed API error code (if available), and the error message.
type APIException struct {
	HTTPStatusCode int
	APIError       *APIError
	ErrorMessage   string
}

func (e *APIException) Error() string {
	if e.APIError != nil {
		return fmt.Sprintf("API error: %v (code: %d, message: %s)", *e.APIError, e.HTTPStatusCode, e.ErrorMessage)
	}
	return fmt.Sprintf("HTTP error: %d (message: %s)", e.HTTPStatusCode, e.ErrorMessage)
}

// APIClient is a client for interacting with the App Store Server API.
// It handles authentication via JWT tokens and provides methods for all API endpoints.
type APIClient struct {
	signingKey  *ecdsa.PrivateKey
	keyID       string
	issuerID    string
	bundleID    string
	environment Environment
	baseURL     string
	httpClient  HTTPClient
}

// NewAPIClient creates a new API client with default HTTP client settings.
// The signingKey should be a PEM-encoded PKCS#8 ECDSA private key.
func NewAPIClient(signingKey []byte, keyID, issuerID, bundleID string, environment Environment) (*APIClient, error) {
	return NewAPIClientWithHTTPClient(signingKey, keyID, issuerID, bundleID, environment, &http.Client{Timeout: 30 * time.Second})
}

// NewAPIClientWithHTTPClient creates a new API client with a custom HTTP client.
// This allows for custom timeout settings, proxies, or mock HTTP clients for testing.
func NewAPIClientWithHTTPClient(signingKey []byte, keyID, issuerID, bundleID string, environment Environment, httpClient HTTPClient) (*APIClient, error) {
	if environment == ENVIRONMENT_XCODE {
		return nil, errors.New("unsupported environment for an APIClient: Xcode")
	}

	block, _ := pem.Decode(signingKey)
	if block == nil {
		return nil, errors.New("failed to parse PEM block from signing key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not an ECDSA private key")
	}

	var baseURL string
	switch environment {
	case ENVIRONMENT_PRODUCTION:
		baseURL = "https://api.storekit.itunes.apple.com"
	case ENVIRONMENT_SANDBOX:
		baseURL = "https://api.storekit-sandbox.itunes.apple.com"
	case ENVIRONMENT_LOCAL_TESTING:
		baseURL = "https://local-testing-base-url"
	default:
		return nil, fmt.Errorf("invalid environment: %v", environment)
	}

	return &APIClient{
		signingKey:  privateKey,
		keyID:       keyID,
		issuerID:    issuerID,
		bundleID:    bundleID,
		environment: environment,
		baseURL:     baseURL,
		httpClient:  httpClient,
	}, nil
}

func (c *APIClient) generateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"bid": c.bundleID,
		"iss": c.issuerID,
		"aud": "appstoreconnect-v1",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	})
	token.Header["kid"] = c.keyID

	return token.SignedString(c.signingKey)
}

func (c *APIClient) makeRequest(method, path string, queryParams url.Values, body, destination any) error {
	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	token, err := c.generateToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "app-store-server-library/go/"+Version())
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	if destination == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(destination)
}

func (c *APIClient) handleErrorResponse(resp *http.Response) error {
	var errorResp struct {
		ErrorCode    int32  `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		return &APIException{
			HTTPStatusCode: resp.StatusCode,
			ErrorMessage:   fmt.Sprintf("failed to decode error response: %v", err),
		}
	}

	apiErr := APIError(errorResp.ErrorCode)
	return &APIException{
		HTTPStatusCode: resp.StatusCode,
		APIError:       &apiErr,
		ErrorMessage:   errorResp.ErrorMessage,
	}
}

// GetTransactionHistory gets a customer's in-app purchase transaction history for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (c *APIClient) GetTransactionHistory(transactionID string, queryParams url.Values, revision string, version GetTransactionHistoryVersion) (*HistoryResponse, error) {
	if revision != "" {
		queryParams.Set("revision", revision)
	}
	if version == "" {
		version = GET_TRANSACTION_HISTORY_VERSION_V1
	}
	path := fmt.Sprintf("/inApps/%s/history/%s", version, transactionID)
	var response HistoryResponse
	if err := c.makeRequest("GET", path, queryParams, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAllSubscriptionStatuses gets the statuses for all of a customer's auto-renewable subscriptions in your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
func (c *APIClient) GetAllSubscriptionStatuses(transactionID string, statuses []Status) (*StatusResponse, error) {
	queryParams := url.Values{}
	for _, status := range statuses {
		queryParams.Add("status", fmt.Sprintf("%d", status))
	}
	path := fmt.Sprintf("/inApps/v1/subscriptions/%s", transactionID)
	var response StatusResponse
	if err := c.makeRequest("GET", path, queryParams, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTransactionInfo gets information about a single transaction for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (c *APIClient) GetTransactionInfo(transactionID string) (*TransactionInfoResponse, error) {
	path := fmt.Sprintf("/inApps/v1/transactions/%s", transactionID)
	var response TransactionInfoResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// LookUpOrderID gets a customer's in-app purchases from a receipt using the order ID.
//
// https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id
func (c *APIClient) LookUpOrderID(orderID string) (*OrderLookupResponse, error) {
	path := fmt.Sprintf("/inApps/v1/lookup/%s", orderID)
	var response OrderLookupResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// RequestTestNotification asks App Store Server Notifications to send a test notification to your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/request_a_test_notification
func (c *APIClient) RequestTestNotification() (*SendTestNotificationResponse, error) {
	path := "/inApps/v1/notifications/test"
	var response SendTestNotificationResponse
	if err := c.makeRequest("POST", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// SendConsumptionInformation sends consumption information about an In-App Purchase to the App Store after your server receives a consumption request notification.
//
// https://developer.apple.com/documentation/appstoreserverapi/send-consumption-information
func (c *APIClient) SendConsumptionInformation(transactionID string, consumptionRequest ConsumptionRequest) error {
	path := fmt.Sprintf("/inApps/v2/transactions/consumption/%s", transactionID)
	return c.makeRequest("PUT", path, nil, consumptionRequest, nil)
}

// SetAppAccountToken sets the app account token value for a purchase the customer makes outside your app, or updates its value in an existing transaction.
//
// https://developer.apple.com/documentation/appstoreserverapi/set-app-account-token
func (c *APIClient) SetAppAccountToken(originalTransactionID string, updateAppAccountTokenRequest UpdateAppAccountTokenRequest) error {
	path := fmt.Sprintf("/inApps/v1/transactions/%s/appAccountToken", originalTransactionID)
	return c.makeRequest("PUT", path, nil, updateAppAccountTokenRequest, nil)
}

// UploadImage uploads an image to use for retention messaging.
//
// https://developer.apple.com/documentation/retentionmessaging/upload-image
func (c *APIClient) UploadImage(imageIdentifier string, image []byte) error {
	path := fmt.Sprintf("/inApps/v1/messaging/image/%s", imageIdentifier)
	return c.makeRequestWithBinaryBody("PUT", path, nil, image, "image/png", nil)
}

// DeleteImage deletes a previously uploaded image.
//
// https://developer.apple.com/documentation/retentionmessaging/delete-image
func (c *APIClient) DeleteImage(imageIdentifier string) error {
	path := fmt.Sprintf("/inApps/v1/messaging/image/%s", imageIdentifier)
	return c.makeRequest("DELETE", path, nil, nil, nil)
}

// GetImageList gets the image identifier and state for all uploaded images.
//
// https://developer.apple.com/documentation/retentionmessaging/get-image-list
func (c *APIClient) GetImageList() (*GetImageListResponse, error) {
	path := "/inApps/v1/messaging/image/list"
	var response GetImageListResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// UploadMessage uploads a message to use for retention messaging.
//
// https://developer.apple.com/documentation/retentionmessaging/upload-message
func (c *APIClient) UploadMessage(messageIdentifier string, uploadMessageRequestBody UploadMessageRequestBody) error {
	path := fmt.Sprintf("/inApps/v1/messaging/message/%s", messageIdentifier)
	return c.makeRequest("PUT", path, nil, uploadMessageRequestBody, nil)
}

// DeleteMessage deletes a previously uploaded message.
//
// https://developer.apple.com/documentation/retentionmessaging/delete-message
func (c *APIClient) DeleteMessage(messageIdentifier string) error {
	path := fmt.Sprintf("/inApps/v1/messaging/message/%s", messageIdentifier)
	return c.makeRequest("DELETE", path, nil, nil, nil)
}

// GetMessageList gets the message identifier and state of all uploaded messages.
//
// https://developer.apple.com/documentation/retentionmessaging/get-message-list
func (c *APIClient) GetMessageList() (*GetMessageListResponse, error) {
	path := "/inApps/v1/messaging/message/list"
	var response GetMessageListResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ConfigureDefaultMessage configures a default message for a specific product in a specific locale.
//
// https://developer.apple.com/documentation/retentionmessaging/configure-default-message
func (c *APIClient) ConfigureDefaultMessage(productID, locale string, defaultConfigurationRequest DefaultConfigurationRequest) error {
	path := fmt.Sprintf("/inApps/v1/messaging/default/%s/%s", productID, locale)
	return c.makeRequest("PUT", path, nil, defaultConfigurationRequest, nil)
}

// DeleteDefaultMessage deletes a default message for a product in a locale.
//
// https://developer.apple.com/documentation/retentionmessaging/delete-default-message
func (c *APIClient) DeleteDefaultMessage(productID, locale string) error {
	path := fmt.Sprintf("/inApps/v1/messaging/default/%s/%s", productID, locale)
	return c.makeRequest("DELETE", path, nil, nil, nil)
}

// GetAppTransactionInfo gets a customer's app transaction information for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/get-app-transaction-info
func (c *APIClient) GetAppTransactionInfo(transactionID string) (*AppTransactionInfoResponse, error) {
	path := fmt.Sprintf("/inApps/v1/transactions/appTransactions/%s", transactionID)
	var response AppTransactionInfoResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *APIClient) makeRequestWithBinaryBody(method, path string, queryParams url.Values, body []byte, contentType string, destination any) error {
	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	token, err := c.generateToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "app-store-server-library/go/"+Version())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	if destination == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(destination)
}

// ExtendRenewalDateForAllActiveSubscribers uses a subscription's product identifier to extend the renewal date for all of its eligible active subscribers.
//
// https://developer.apple.com/documentation/appstoreserverapi/extend_subscription_renewal_dates_for_all_active_subscribers
func (c *APIClient) ExtendRenewalDateForAllActiveSubscribers(massExtendRenewalDateRequest MassExtendRenewalDateRequest) (*MassExtendRenewalDateResponse, error) {
	path := "/inApps/v1/subscriptions/extend/mass"
	var response MassExtendRenewalDateResponse
	if err := c.makeRequest("POST", path, nil, massExtendRenewalDateRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ExtendSubscriptionRenewalDate extends the renewal date of a customer's active subscription using the original transaction identifier.
//
// https://developer.apple.com/documentation/appstoreserverapi/extend_a_subscription_renewal_date
func (c *APIClient) ExtendSubscriptionRenewalDate(originalTransactionID string, extendRenewalDateRequest ExtendRenewalDateRequest) (*ExtendRenewalDateResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/%s", originalTransactionID)
	var response ExtendRenewalDateResponse
	if err := c.makeRequest("PUT", path, nil, extendRenewalDateRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRefundHistory gets a paginated list of all of a customer's refunded in-app purchases for your app.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (c *APIClient) GetRefundHistory(transactionID, revision string) (*RefundHistoryResponse, error) {
	queryParams := url.Values{}
	if revision != "" {
		queryParams.Set("revision", revision)
	}
	path := fmt.Sprintf("/inApps/v2/refund/lookup/%s", transactionID)
	var response RefundHistoryResponse
	if err := c.makeRequest("GET", path, queryParams, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetStatusOfSubscriptionRenewalDateExtensions checks whether a renewal date extension request completed, and provides the final count of successful or failed extensions.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_status_of_subscription_renewal_date_extensions
func (c *APIClient) GetStatusOfSubscriptionRenewalDateExtensions(requestIdentifier, productID string) (*MassExtendRenewalDateStatusResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/mass/%s/%s", requestIdentifier, productID)
	var response MassExtendRenewalDateStatusResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTestNotificationStatus checks the status of the test App Store server notification sent to your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_test_notification_status
func (c *APIClient) GetTestNotificationStatus(testNotificationToken string) (*CheckTestNotificationResponse, error) {
	path := fmt.Sprintf("/inApps/v1/notifications/test/%s", testNotificationToken)
	var response CheckTestNotificationResponse
	if err := c.makeRequest("GET", path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetNotificationHistory gets a list of notifications that the App Store server attempted to send to your server.
//
// https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
func (c *APIClient) GetNotificationHistory(paginationToken string, notificationHistoryRequest NotificationHistoryRequest) (*NotificationHistoryResponse, error) {
	queryParams := url.Values{}
	if paginationToken != "" {
		queryParams.Set("paginationToken", paginationToken)
	}
	path := "/inApps/v1/notifications/history"
	var response NotificationHistoryResponse
	if err := c.makeRequest("POST", path, queryParams, notificationHistoryRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
