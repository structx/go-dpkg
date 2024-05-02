// Package encode util helper functions for encoding values
package encode

// MaskFromPrefix created zero filled mask shifting prefix
func MaskFromPrefix(nodeID [28]byte, level int) [28]byte {

	var mask [28]byte

	// initialize mask zero filled
	for i := range mask {
		mask[i] = 0x00
	}

	// calculate prefix length
	prefixLength := uint(8) * uint(level)

	// create prefix mask by shifting left
	for i := range mask {
		// invert the index to create a left shift based on prefix length
		shift := uint(len(mask)) - uint(i) - 1
		if shift < prefixLength {
			mask[i] = 0xFF // set bits to 1 based on inverted index
		} else {
			mask[i] = 0x00 // set remaining bits to 0
		}
	}

	result := [28]byte{}
	// perform bitwise AND to extract prefix
	for i := range result {
		result[i] = nodeID[i] & mask[i]
	}

	return result
}
