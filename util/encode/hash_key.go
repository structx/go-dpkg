package encode

import "golang.org/x/crypto/sha3"

// HashKey use sha3 to hash key
func HashKey(key []byte) [28]byte {

	h := sha3.New224()
	h.Write(key)
	hash := h.Sum(nil)

	var result [28]byte
	copy(result[:], hash[:28])

	return result
}
