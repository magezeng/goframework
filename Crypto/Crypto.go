package Crypto

import (
	"crypto/rand"
	"io"
)

/*
对称加密：
AES-128(已实现)
AES-192(已实现)
AES-256(已实现)
非对称加密：
RSA
哈希算法：
MD5
HMAC
SHA
*/

const (
	// message 12个字节 目前用明文，也可以从加密串里面来取得
	DEFAULT_SECRET_MESSAGE = "secretmessag"
	// key 16个字节 目前用明文，也可以从加密串里面来取得
	DEFAULT_SECRETKEY = "Tipu!@#123456789"
)

// generateNonce 获得指定长度的随机字节数组
func generateNonce(length uint64) (nonce []byte, err error) {
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce = make([]byte, length)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}
	return
}
