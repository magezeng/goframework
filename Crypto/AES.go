package Crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type AES struct {
}

// Encrypt AES用gcm模式加密
func (AES) Encrypt(originData string) (result []byte, err error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	result = gcm.Seal(nil, secretMessage, []byte(originData), nil)
	return
}

// Encrypt AES用gcm模式解密
func (AES) Decrypt(encryptedData string) (result []byte, err error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	result, err = gcm.Open(nil, secretMessage, []byte(encryptedData), nil)
	if err != nil {
		return
	}
	return
}

//AES key必须是16字节 24字节 32字节中的一种
func checkValidAESKey(strKey string) []byte {
	keyLen := len(strKey)
	arrKey := []byte(strKey)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	if keyLen >= 16 {
		//取前16个字节
		return arrKey[:16]
	}
	//补齐16个字节
	tmp := make([]byte, 16)
	for i := 0; i < 16; i++ {
		if i < keyLen {
			tmp[i] = arrKey[i]
		} else {
			tmp[i] = '0'
		}
	}
	return tmp
}
