package appstore

import (
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"regexp"
)

const (
	pkcs7Oid                      = "1.2.840.113549.1.7.2"
	inAppArray                    = 17
	transactionIdentifier         = 1703
	originalTransactionIdentifier = 1705
)

// ReceiptUtility provides methods for extracting transaction IDs from receipts.
type ReceiptUtility struct {
}

// NewReceiptUtility creates a new ReceiptUtility instance.
func NewReceiptUtility() *ReceiptUtility {
	return &ReceiptUtility{}
}

// ExtractTransactionIdFromAppReceipt extracts a transaction id from an encoded App Receipt.
// Returns an error if the receipt does not match the expected format.
// NO validation is performed on the receipt, and any data returned should only be used to call the App Store Server API.
//
// See https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (u *ReceiptUtility) ExtractTransactionIdFromAppReceipt(appReceipt string) (*string, error) {
	data, err := base64.StdEncoding.DecodeString(appReceipt)
	if err != nil {
		return nil, err
	}

	decoder := newIndefiniteFormAwareDecoder(data)
	tag, err := decoder.peek()
	if err != nil {
		return nil, err
	}
	if !tag.isConstructed || tag.number != uint32(asn1.TagSequence) {
		return nil, errors.New("invalid receipt format")
	}
	if err := decoder.enter(); err != nil {
		return nil, err
	}

	// PKCS#7 object
	tag, value, err := decoder.read()
	if err != nil {
		return nil, err
	}
	if tag.isConstructed || tag.number != uint32(asn1.TagOID) || value.(asn1.ObjectIdentifier).String() != pkcs7Oid {
		return nil, errors.New("invalid receipt format")
	}

	// This is the PKCS#7 format, work our way into the inner content
	for _, shouldEnter := range []bool{true, true, false, false, true, false, true} {
		var err error
		if shouldEnter {
			err = decoder.enter()
		} else {
			_, _, err = decoder.read()
		}
		if err != nil {
			return nil, err
		}
	}

	tag, value, err = decoder.read()
	if err != nil {
		return nil, err
	}

	// Xcode uses nested OctetStrings, we extract the inner string in this case
	if tag.isConstructed && tag.number == uint32(asn1.TagOctetString) {
		innerDecoder := newAsn1Decoder(value.([]byte))
		tag, value, err = innerDecoder.read()
		if err != nil {
			return nil, err
		}
	}

	if tag.number != uint32(asn1.TagOctetString) {
		return nil, errors.New("invalid receipt format")
	}

	decoder = newAsn1Decoder(value.([]byte))
	tag, err = decoder.peek()
	if err != nil {
		return nil, err
	}
	if !tag.isConstructed || tag.number != uint32(asn1.TagSet) {
		return nil, errors.New("invalid receipt format")
	}
	if err := decoder.enter(); err != nil {
		return nil, err
	}

	// We are in the top-level sequence, work our way to the array of in-apps
	for !decoder.eof() {
		if err := decoder.enter(); err != nil {
			break
		}
		tag, value, err = decoder.read()
		if err == nil && !tag.isConstructed && tag.number == uint32(asn1.TagInteger) && value.(int64) == inAppArray {
			if _, _, err = decoder.read(); err != nil {
				return nil, err
			}
			tag, value, err = decoder.read()
			if err == nil && !tag.isConstructed && tag.number == uint32(asn1.TagOctetString) {
				inAppDecoder := newAsn1Decoder(value.([]byte))
				if err := inAppDecoder.enter(); err != nil {
					return nil, err
				}
				// In-app array
				for !inAppDecoder.eof() {
					if err := inAppDecoder.enter(); err != nil {
						break
					}
					tag, value, err = inAppDecoder.read()
					if err == nil && !tag.isConstructed && tag.number == uint32(asn1.TagInteger) && (value.(int64) == transactionIdentifier || value.(int64) == originalTransactionIdentifier) {
						if _, _, err = inAppDecoder.read(); err != nil {
							return nil, err
						}
						tag, value, err = inAppDecoder.read()
						if err == nil {
							singletonDecoder := newAsn1Decoder(value.([]byte))
							_, singletonValue, err := singletonDecoder.read()
							if err == nil {
								transactionId := singletonValue.(string)
								return &transactionId, nil
							}
						}
					}
					inAppDecoder.leave()
				}
			}
		}
		decoder.leave()
	}

	return nil, nil
}

// ExtractTransactionIdFromTransactionReceipt extracts a transaction id from an encoded transactional receipt.
// Returns an error if the receipt does not match the expected format.
// NO validation is performed on the receipt, and any data returned should only be used to call the App Store Server API.
func (u *ReceiptUtility) ExtractTransactionIdFromTransactionReceipt(transactionReceipt string) (*string, error) {
	decodedTopLevel, err := base64.StdEncoding.DecodeString(transactionReceipt)
	if err != nil {
		return nil, err
	}
	topLevelStr := string(decodedTopLevel)

	purchaseInfoRegex := regexp.MustCompile(`"purchase-info"\s+=\s+"([a-zA-Z0-9+/=]+)";`)
	matchingResult := purchaseInfoRegex.FindStringSubmatch(topLevelStr)
	if len(matchingResult) > 1 {
		decodedInnerLevel, err := base64.StdEncoding.DecodeString(matchingResult[1])
		if err != nil {
			return nil, err
		}
		innerLevelStr := string(decodedInnerLevel)

		transactionIdRegex := regexp.MustCompile(`"transaction-id"\s+=\s+"([a-zA-Z0-9+/=]+)";`)
		innerMatchingResult := transactionIdRegex.FindStringSubmatch(innerLevelStr)
		if len(innerMatchingResult) > 1 {
			res := innerMatchingResult[1]
			return &res, nil
		}
	}
	return nil, nil
}

type asn1Tag struct {
	isConstructed bool
	number        uint32
}

type decoderState struct {
	data   []byte
	offset int
}

type asn1Decoder struct {
	stack                 []*decoderState
	isIndefiniteFormAware bool
}

func newAsn1Decoder(data []byte) *asn1Decoder {
	return &asn1Decoder{
		stack:                 []*decoderState{{data: data, offset: 0}},
		isIndefiniteFormAware: false,
	}
}

func newIndefiniteFormAwareDecoder(data []byte) *asn1Decoder {
	return &asn1Decoder{
		stack:                 []*decoderState{{data: data, offset: 0}},
		isIndefiniteFormAware: true,
	}
}

func (d *asn1Decoder) currentState() *decoderState {
	if len(d.stack) == 0 {
		return nil
	}
	return d.stack[len(d.stack)-1]
}

func (d *asn1Decoder) eof() bool {
	state := d.currentState()
	return state == nil || state.offset >= len(state.data)
}

func (d *asn1Decoder) peek() (asn1Tag, error) {
	state := d.currentState()
	if state == nil || state.offset >= len(state.data) {
		return asn1Tag{}, errors.New("EOF")
	}
	tagByte := state.data[state.offset]
	return asn1Tag{
		isConstructed: (tagByte & 0x20) != 0,
		number:        uint32(tagByte & 0x1F),
	}, nil
}

func (d *asn1Decoder) read() (asn1Tag, any, error) {
	state := d.currentState()
	tag, err := d.peek()
	if err != nil {
		return asn1Tag{}, nil, err
	}
	state.offset++

	length, err := d.readLength()
	if err != nil {
		return asn1Tag{}, nil, err
	}

	if state.offset+length > len(state.data) {
		return asn1Tag{}, nil, errors.New("offset out of range")
	}

	valueData := state.data[state.offset : state.offset+length]
	state.offset += length

	var value any
	switch tag.number {
	case uint32(asn1.TagInteger):
		var i int64
		_, err = asn1.Unmarshal(append([]byte{byte(tag.number & 0x1F)}, append(encodeLength(length), valueData...)...), &i)
		value = i
	case uint32(asn1.TagOctetString):
		value = valueData
	case uint32(asn1.TagOID):
		var oid asn1.ObjectIdentifier
		_, err = asn1.Unmarshal(append([]byte{byte(tag.number & 0x1F)}, append(encodeLength(length), valueData...)...), &oid)
		value = oid
	case uint32(asn1.TagUTF8String), uint32(asn1.TagPrintableString), uint32(asn1.TagIA5String):
		value = string(valueData)
	default:
		value = valueData
	}

	return tag, value, err
}

func (d *asn1Decoder) readLength() (int, error) {
	state := d.currentState()
	if state == nil || state.offset >= len(state.data) {
		return 0, errors.New("EOF")
	}
	b := state.data[state.offset]
	if d.isIndefiniteFormAware && b == 0x80 {
		state.offset++
		return len(state.data) - state.offset, nil
	}
	state.offset++
	if b < 0x80 {
		return int(b), nil
	}
	numBytes := int(b & 0x7F)
	if numBytes == 0 {
		return 0, errors.New("indefinite length not supported")
	}
	length := 0
	for range numBytes {
		if state.offset >= len(state.data) {
			return 0, errors.New("EOF")
		}
		length = (length << 8) | int(state.data[state.offset])
		state.offset++
	}
	return length, nil
}

func (d *asn1Decoder) enter() error {
	state := d.currentState()
	_, err := d.peek()
	if err != nil {
		return err
	}
	state.offset++

	length, err := d.readLength()
	if err != nil {
		return err
	}

	if state.offset+length > len(state.data) {
		return errors.New("offset out of range")
	}

	content := state.data[state.offset : state.offset+length]
	state.offset += length

	d.stack = append(d.stack, &decoderState{data: content, offset: 0})
	return nil
}

func (d *asn1Decoder) leave() {
	if len(d.stack) > 1 {
		d.stack = d.stack[:len(d.stack)-1]
	}
}

func encodeLength(length int) []byte {
	if length < 0x80 {
		return []byte{byte(length)}
	}
	var b []byte
	for length > 0 {
		b = append([]byte{byte(length & 0xFF)}, b...)
		length >>= 8
	}
	return append([]byte{byte(len(b) | 0x80)}, b...)
}
