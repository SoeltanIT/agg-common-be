package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SignaturePayload represents the internal data structure encoded in the signature
type SignaturePayload struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
	Hash      string `json:"hash"`
}

// ValidateSignature validates signature and extracts its data
func ValidateSignature(signatureB64 string) error {
	jsonData, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("invalid signature format: cannot decode base64: %w", err)
	}

	var payload SignaturePayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return fmt.Errorf("invalid signature format: cannot parse JSON: %w", err)
	}

	if payload.Token == "" || payload.Timestamp == 0 || payload.Hash == "" {
		return errors.New("invalid signature payload: missing fields")
	}

	now := time.Now().Unix()
	expiration := int64(5 * 60)
	if now-payload.Timestamp > expiration {
		return fmt.Errorf("signature expired: %d seconds old", now-payload.Timestamp)
	}
	if payload.Timestamp > now+60 {
		return errors.New("invalid timestamp: timestamp from future")
	}

	data := fmt.Sprintf("%s%d", payload.Token, payload.Timestamp)
	h := hmac.New(sha256.New, []byte(payload.Token))
	h.Write([]byte(data))
	expectedHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(payload.Hash), []byte(expectedHash)) {
		return errors.New("invalid signature: hash mismatch")
	}

	return nil
}
