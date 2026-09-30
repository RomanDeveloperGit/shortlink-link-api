package shortcode

import (
	"crypto/rand"
	"math/big"
)

var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func GenerateShortCode(size int) (string, error) {
	b := make([]byte, size)
	maxIdx := big.NewInt(int64(len(chars)))

	for i := range b {
		n, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", err
		}

		b[i] = chars[n.Int64()]
	}

	return string(b), nil
}
