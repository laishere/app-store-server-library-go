package appstore

import (
	"encoding/json"
	"testing"
)

func TestRealtimeResponseBodyWithMessage(t *testing.T) {
	messageId := "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f67890"
	message := &Message{
		MessageIdentifier: &messageId,
	}
	responseBody := &RealtimeResponseBody{
		Message: message,
	}

	jsonData, err := json.Marshal(responseBody)
	assertNoError(t, err, "Failed to marshal response")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assertNoError(t, err, "Failed to unmarshal to dict")

	assertNotNil(t, jsonDict["message"], "message in jsonDict")
	msg := jsonDict["message"].(map[string]any)
	assertEqual(t, messageId, msg["messageIdentifier"], "messageIdentifier")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assertNoError(t, err, "Failed to deserialize")

	assertNotNil(t, deserialized.Message, "Message")
	assertEqual(t, messageId, *deserialized.Message.MessageIdentifier, "MessageIdentifier")
}

func TestRealtimeResponseBodyWithAlternateProduct(t *testing.T) {
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
	assertNoError(t, err, "Failed to marshal")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assertNoError(t, err, "Failed to unmarshal to dict")

	assertNotNil(t, jsonDict["alternateProduct"], "alternateProduct in jsonDict")
	alt := jsonDict["alternateProduct"].(map[string]any)
	assertEqual(t, messageId, alt["messageIdentifier"], "messageIdentifier")
	assertEqual(t, productId, alt["productId"], "productId")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assertNoError(t, err, "Failed to deserialize")

	assertNotNil(t, deserialized.AlternateProduct, "AlternateProduct")
	assertEqual(t, messageId, *deserialized.AlternateProduct.MessageIdentifier, "MessageIdentifier")
	assertEqual(t, productId, *deserialized.AlternateProduct.ProductId, "ProductId")
}

func TestRealtimeResponseBodyWithPromotionalOfferV2(t *testing.T) {
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
	assertNoError(t, err, "Failed to marshal")

	var jsonDict map[string]any
	err = json.Unmarshal(jsonData, &jsonDict)
	assertNoError(t, err, "Failed to unmarshal to dict")

	assertNotNil(t, jsonDict["promotionalOffer"], "promotionalOffer in jsonDict")
	offer := jsonDict["promotionalOffer"].(map[string]any)
	assertEqual(t, messageId, offer["messageIdentifier"], "messageIdentifier")
	assertEqual(t, signatureV2, offer["promotionalOfferSignatureV2"], "promotionalOfferSignatureV2")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assertNoError(t, err, "Failed to deserialize")

	assertNotNil(t, deserialized.PromotionalOffer, "PromotionalOffer")
	assertEqual(t, messageId, *deserialized.PromotionalOffer.MessageIdentifier, "MessageIdentifier")
	assertEqual(t, signatureV2, *deserialized.PromotionalOffer.PromotionalOfferSignatureV2, "PromotionalOfferSignatureV2")
}

func TestRealtimeResponseBodyWithPromotionalOfferV1(t *testing.T) {
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
	assertNoError(t, err, "Failed to marshal")

	var deserialized RealtimeResponseBody
	err = json.Unmarshal(jsonData, &deserialized)
	assertNoError(t, err, "Failed to deserialize")

	assertNotNil(t, deserialized.PromotionalOffer, "PromotionalOffer")
	v1 := deserialized.PromotionalOffer.PromotionalOfferSignatureV1
	assertEqual(t, productId, *v1.ProductId, "ProductId")
	assertEqual(t, offerId, *v1.OfferIdentifier, "OfferIdentifier")
	assertEqual(t, nonce, *v1.Nonce, "Nonce")
	assertEqual(t, timestamp, *v1.Timestamp, "Timestamp")
	assertEqual(t, keyId, *v1.KeyId, "KeyId")
	assertEqual(t, appAccountToken, *v1.AppAccountToken, "AppAccountToken")
	assertEqual(t, encodedSignature, *v1.EncodedSignature, "EncodedSignature")
}
