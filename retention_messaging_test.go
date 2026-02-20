package appstore

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRealtimeResponseBodyWithMessage(t *testing.T) {
	assert := assert.New(t)
	messageId := "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f67890"
	message := &Message{
		MessageIdentifier: &messageId,
	}
	responseBody := &RealtimeResponseBody{
		Message: message,
	}

	jsonData, err := json.Marshal(responseBody)
	assert.NoError(err, "Failed to marshal response")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assert.NoError(err, "Failed to unmarshal to dict")

	assert.NotNil(jsonDict["message"], "message in jsonDict")
	msg := jsonDict["message"].(map[string]any)
	assert.Equal(messageId, msg["messageIdentifier"], "messageIdentifier")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assert.NoError(err, "Failed to deserialize")

	assert.NotNil(deserialized.Message, "Message")
	assert.Equal(messageId, *deserialized.Message.MessageIdentifier, "MessageIdentifier")
}

func TestRealtimeResponseBodyWithAlternateProduct(t *testing.T) {
	assert := assert.New(t)
	messageId := "b2c3d4e5-f6a7-8901-b2c3-d4e5f6a78901"
	productId := "com.example.alternate.product"
	alternateProduct := &AlternateProduct{
		MessageIdentifier: &messageId,
		ProductId:         &productId,
	}
	responseBody := &RealtimeResponseBody{
		AlternateProduct: alternateProduct,
	}

	jsonData, err := json.Marshal(responseBody)
	assert.NoError(err, "Failed to marshal")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assert.NoError(err, "Failed to unmarshal to dict")

	assert.NotNil(jsonDict["alternateProduct"], "alternateProduct in jsonDict")
	alt := jsonDict["alternateProduct"].(map[string]any)
	assert.Equal(messageId, alt["messageIdentifier"], "messageIdentifier")
	assert.Equal(productId, alt["productId"], "productId")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assert.NoError(err, "Failed to deserialize")

	assert.NotNil(deserialized.AlternateProduct, "AlternateProduct")
	assert.Equal(messageId, *deserialized.AlternateProduct.MessageIdentifier, "MessageIdentifier")
	assert.Equal(productId, *deserialized.AlternateProduct.ProductId, "ProductId")
}

func TestRealtimeResponseBodyWithPromotionalOfferV2(t *testing.T) {
	assert := assert.New(t)
	messageId := "c3d4e5f6-a789-0123-c3d4-e5f6a7890123"
	signatureV2 := "signature2"
	promotionalOffer := &PromotionalOffer{
		MessageIdentifier:           &messageId,
		PromotionalOfferSignatureV2: &signatureV2,
	}
	responseBody := &RealtimeResponseBody{
		PromotionalOffer: promotionalOffer,
	}

	jsonData, err := json.Marshal(responseBody)
	assert.NoError(err, "Failed to marshal")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assert.NoError(err, "Failed to unmarshal to dict")

	assert.NotNil(jsonDict["promotionalOffer"], "promotionalOffer in jsonDict")
	offer := jsonDict["promotionalOffer"].(map[string]any)
	assert.Equal(messageId, offer["messageIdentifier"], "messageIdentifier")
	assert.Equal(signatureV2, offer["promotionalOfferSignatureV2"], "promotionalOfferSignatureV2")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assert.NoError(err, "Failed to deserialize")

	assert.NotNil(deserialized.PromotionalOffer, "PromotionalOffer")
	assert.Equal(messageId, *deserialized.PromotionalOffer.MessageIdentifier, "MessageIdentifier")
	assert.Equal(signatureV2, *deserialized.PromotionalOffer.PromotionalOfferSignatureV2, "PromotionalOfferSignatureV2")
}

func TestRealtimeResponseBodyWithPromotionalOfferV1(t *testing.T) {
	assert := assert.New(t)
	messageId := "d4e5f6a7-8901-2345-d4e5-f6a789012345"
	nonce := "e5f6a789-0123-4567-e5f6-a78901234567"
	appAccountToken := "f6a78901-2345-6789-f6a7-890123456789"
	timestamp := Timestamp(1698148900000)
	productId := "com.example.product"
	keyId := "keyId123"
	offerId := "offer123"
	encodedSignature := "base64encodedSignature"

	signatureV1 := &PromotionalOfferSignatureV1{
		EncodedSignature: &encodedSignature,
		ProductId:        &productId,
		Nonce:            &nonce,
		Timestamp:        &timestamp,
		KeyId:            &keyId,
		OfferIdentifier:  &offerId,
		AppAccountToken:  &appAccountToken,
	}

	promotionalOffer := &PromotionalOffer{
		MessageIdentifier:           &messageId,
		PromotionalOfferSignatureV1: signatureV1,
	}
	responseBody := &RealtimeResponseBody{
		PromotionalOffer: promotionalOffer,
	}

	jsonData, err := json.Marshal(responseBody)
	assert.NoError(err, "Failed to marshal")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assert.NoError(err, "Failed to deserialize")

	assert.NotNil(deserialized.PromotionalOffer, "PromotionalOffer")
	v1 := deserialized.PromotionalOffer.PromotionalOfferSignatureV1
	assert.Equal(productId, *v1.ProductId, "ProductId")
	assert.Equal(offerId, *v1.OfferIdentifier, "OfferIdentifier")
	assert.Equal(nonce, *v1.Nonce, "Nonce")
	assert.Equal(timestamp, *v1.Timestamp, "Timestamp")
	assert.Equal(keyId, *v1.KeyId, "KeyId")
	assert.Equal(appAccountToken, *v1.AppAccountToken, "AppAccountToken")
	assert.Equal(encodedSignature, *v1.EncodedSignature, "EncodedSignature")
}
