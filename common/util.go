package common

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	md5ctx := hex.EncodeToString(h.Sum(nil))
	return md5ctx
}
