package utils

import (
	"bytes"
	"encoding/hex"
	"strings"
)

func MergeBytes(arr ...[]byte) (data []byte) {
	pBytes := arr
	l := len(pBytes)
	s := make([][]byte, l)
	for index := 0; index < l; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	data = bytes.Join(s, sep)
	return
}

func CutBytes(d []byte, offset, length int) (data []byte) {
	data = d[offset : offset+length]
	return
}

func BytesToHex(data []byte) (hexStr string) {
	hexStr = strings.ToUpper(hex.EncodeToString(data))
	return
}

func HexToBytes(hexStr string) (data []byte) {
	data, _ = hex.DecodeString(hexStr)
	return
}
