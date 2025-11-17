package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeHash(body []byte, key string) string {

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
