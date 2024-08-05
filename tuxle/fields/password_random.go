package fields

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const SALT_LENGTH = 8

func encodeString(str string) string {
	var sha512Hasher = sha512.New()
	sha512Hasher.Write([]byte(str))
	return hex.EncodeToString(sha512Hasher.Sum(nil))
}

func RandomString(length int) string {
	if length > 128 {
		numCalls := (length + 127) / 128 // Ceiling of (length / 128)

		var builder strings.Builder
		for i := range numCalls {
			builder.WriteString(encodeString(fmt.Sprint(time.Now().UnixNano() + int64(i))))
		}

		return builder.String()[:length]
	}

	return encodeString(fmt.Sprint(time.Now().UnixNano()))[0:length]
}
