package utils

import (
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

func Sha3(str string) string {
	h := sha3.New512()
	h.Write([]byte(str))
	hash := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash[:])

}

func Reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return result
}

func RandomInt(strlen int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	const chars = "1234567890"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}
	conv, _ := strconv.Atoi(string(result))

	return conv
}
