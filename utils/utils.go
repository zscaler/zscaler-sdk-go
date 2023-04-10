package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

func IsTokenExpired(tokenString string) bool {
	// Split the token into three parts: header, payload, and signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return true
	}

	// Decode the payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}

	// Parse the payload as JSON
	var payload map[string]interface{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return true
	}

	// Check the expiration time
	if exp, ok := payload["exp"].(float64); ok {
		// minus 10 seconds to avoid token expired
		exp = exp - 10
		if time.Now().Unix() > int64(exp) {
			return true
		}
	}

	return false
}
