package Crypto

import (
	"crypto/md5"
	"fmt"
)

/*
  md5 sign
*/
func Md5Signer(message string) string {
	data := []byte(message)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}