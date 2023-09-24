package utils

import (
	"encoding/hex"
	"fmt"
)

func Bytes2Hex(data []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(data))
}
