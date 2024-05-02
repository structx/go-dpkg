package encode

// GenerateBucketID
func GenerateBucketID(nodeID [28]byte, level int) [28]byte {

	var mask [28]byte

	// Initialize mask with all 0s
	for i := range mask {
		mask[i] = 0x00
	}

	// Calculate prefix length
	prefixLength := uint(8) * uint(level)

	// Create prefix mask by shifting left
	for i := range mask {
		// Invert the index to create a left shift based on prefix length
		shift := uint(len(mask)) - uint(i) - 1
		if shift < prefixLength {
			mask[i] = 0xFF // Set bits to 1 based on inverted index
		} else {
			mask[i] = 0x00 // Set remaining bits to 0
		}
	}

	result := [28]byte{}
	// Perform bitwise AND to extract prefix
	for i := range result {
		result[i] = nodeID[i] & mask[i]
	}

	return result
}
