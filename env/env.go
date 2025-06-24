package env

import (
	"crypto/sha256"
	"encoding/base32"
	"os"
	"strings"
	"unicode"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func Id(target string) string {
	hash := sha256.New()
	hash.Write([]byte(target))
	hashBytes := hash.Sum(nil)

	encoded := base32.StdEncoding.EncodeToString(hashBytes)

	if len(encoded) > 10 {
		encoded = encoded[:10]
	}

	if !unicode.IsLetter(rune(encoded[0])) {
		replacement := replace(encoded)
		encoded = string(replacement) + encoded[1:]
	}

	return strings.ToLower(encoded)
}

func Get(either string, or ...string) string {
	if value, exists := os.LookupEnv(either); exists {
		return value
	}

	if len(or) > 0 {
		return or[0]
	}

	return ""
}

func replace(s string) byte {
	index := int(s[0]) % len(letters)
	if index < 0 {
		index = -index
	}
	return letters[index]
}
