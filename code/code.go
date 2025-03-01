package code

import "math/rand/v2"

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode() string {
	result := make([]byte, 6)

	for i := range result {
		result[i] = chars[rand.IntN(len(chars))]
	}

	return string(result)
}
