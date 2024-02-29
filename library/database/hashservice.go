package database

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func GetHash(values ...string) string {
	hasher := sha1.New()
	hasher.Write([]byte(strings.Join(values, "")))
	hex_data := hex.EncodeToString(hasher.Sum(nil))
	return hex_data
}
