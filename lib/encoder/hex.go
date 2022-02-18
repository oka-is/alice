package encoder

import "encoding/hex"

func B2Hex(input []byte) string {
	return hex.EncodeToString(input)
}
