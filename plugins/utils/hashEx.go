package utils

import (
	"crypto/hmac"
	"crypto/md5"
)

func GetBytesHmac(d, key []byte) (data []byte) {
	h := hmac.New(md5.New, key)
	h.Write(d)
	data = h.Sum(nil)
	return
}
