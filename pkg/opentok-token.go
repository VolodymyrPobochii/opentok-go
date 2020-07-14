package pkg

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

const TokenSentinel = "T1=="

// @typedef {Object} TokenData
// @property {string} [session_id] An OpenTok Session ID
// @property {number} [create_time] Creation time of token as unix timestamp (Default: now)
// @property {number} [expire_time] Expiration time of token as unix timestamp (Default: one day
// from now)
// @property {number} [nonce] Arbitrary number used only once in a cryptographic communication
// (Default: unique random number)
// @property {string} [role='publisher'] "publisher" or "subscriber" "moderator"
// @property {string} [connection_data] Arbitrary data to be made available in clients on the OpenTok Connection

// Encodes data for use as a token that can be used as the X-TB-TOKEN-AUTH header value in OpenTok REST APIs
//
// @exports opentok-token
//
// @param {TokenData} tokenData
// @param {string} apiKey An OpenTok API Key
// @param {string} apiSecret An OpenTok API Secret
//
// @returns {string} token
func EncodeToken(tokenData map[string]interface{}, apiKey, apiSecret string) (string, error) {
	tokenData = Clone(tokenData)
	now := time.Now()
	nonce, err := Nonce(0)()
	if err != nil {
		return "", err
	}
	tokenData = Defaults(tokenData, map[string]interface{}{
		"create_time": now.UnixNano() / int64(time.Millisecond),
		"expire_time": now.Add(24*time.Hour).UnixNano() / int64(time.Millisecond),
		"nonce":       nonce,
		"role":        "publisher",
	})
	dataString := QueryString(tokenData)
	sig, err := signString(dataString, apiSecret)
	if err != nil {
		return "", err
	}
	decoded := fmt.Sprintf("partner_id=%s&sig=%s:%s", apiKey, sig, dataString)
	return fmt.Sprintf("%s%s", TokenSentinel, base64.StdEncoding.EncodeToString([]byte(decoded))), nil
}

// Creates an HMAC-SHA1 signature of unsigned data using the key
//
// @private
//
// @param {string} unsigned Data to be signed
// @param {string} key Key to sign data with
//
// @returns {string} signature
func signString(unsigned, secret string) (string, error) {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha1.New, []byte(secret))
	// Write Data to it
	_, err := h.Write([]byte(unsigned))
	if err != nil {
		return "", err
	}
	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	return sha, nil
}
