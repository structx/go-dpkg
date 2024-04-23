package structs

import (
	"golang.org/x/crypto/sha3"
)

func hashKey(key []byte) [28]byte {
	h := sha3.New224()
	h.Write(key)
	hash := h.Sum(nil)

	var result [28]byte
	copy(result[:], hash[:28])

	return result
}

func xor(a, b [28]byte) [28]byte {
	var result [28]byte
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func compareDistances(a, b [28]byte) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
