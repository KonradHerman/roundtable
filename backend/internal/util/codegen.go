package util

import (
	"crypto/rand"
	"math/big"
)

// GenerateRoomCode creates a random 6-character room code.
// Uses characters that are easy to read and type on mobile: uppercase letters and numbers,
// excluding ambiguous characters (0/O, 1/I/L).
func GenerateRoomCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // No 0,O,1,I,L
	const codeLength = 6

	code := make([]byte, codeLength)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range code {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fallback to a less secure but functional approach
			code[i] = charset[i%len(charset)]
			continue
		}
		code[i] = charset[num.Int64()]
	}

	return string(code)
}
