package engine

import (
	"math/rand"
)

const (
	CENTISECONDS_PER_SECOND = 100
)

/*
	Converts a decimal value to a little-endian byte array. This is 
	useful for converting values to the [lo, hi] or [lo, mid, hi] 
	format used in PM5 characteristics.
	E.g. decimalToHexBytes(12000, 2) returns [224, 46]
	because 224 + (46 * 256) = 12000
*/
func decimalToHexBytes(value int, len int) []byte {	
	bytes := make([]byte, len, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(value % 256)
		value = int(value / 256)	
	}
	return bytes
}

/*
	Generates the lo and hi bytes for a realistic "Split/Int Avg Pace" value
	For now "realistic" is defined as a random value between 120s and 130s,
	corresponding to a 2:00-2:10 split. 
*/
func generateSplitBytes() []byte {
	splitCentiseconds := (120 + rand.Intn(10)) * CENTISECONDS_PER_SECOND
	return decimalToHexBytes(splitCentiseconds, 2)
}

/*
	Generates a 20 byte characteristic
	Currently only the Split/Int bytes are non-zero.
	More coming soon!
*/
func GenerateAdditionalStatus2Char() []byte {
	bytes := make([]byte, 20, 20)
	copy(bytes[8:], generateSplitBytes()) // split/int lo and hi start at byte 8
	return bytes
}
