package Crypto

import (
	"crypto/rand"
	"io"
)

/*
对称加密：
AES-128（已实现）
非对称加密：
RSA
哈希算法：
MD5
HMAC
SHA
*/

var (
	// message 12个字节 目前用明文，也可以从加密串里面来取得
	secretMessage = []byte("secretmessag")
	// key 16个字节 目前用明文，也可以从加密串里面来取得
	secretKey = []byte("Tipu!@#123456789")
)

// GenerateNonce 获得指定长度的随机字节数组
func GenerateNonce(length uint64) (nonce []byte, err error) {
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce = make([]byte, length)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}
	return
}