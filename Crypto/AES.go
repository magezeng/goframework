package Crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type AES struct {
}

// Encrypt AES用gcm模式默认加密
func (a AES) Encrypt(originData string) (result []byte, err error) {
	return a.EncryptWithCustomData(originData, DEFAULT_SECRETKEY, DEFAULT_SECRET_MESSAGE)
}

func (a AES) EncryptWithCustomData(originData string, secretKey string, nonce string) (result []byte, err error) {
	// 获得有效的aes密钥
	key := getValidAESKey(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	// 没有nonce时生成新的nonce
	var nonceParam []byte
	if nonce == "" {
		nonceParam, err = generateNonce(uint64(gcm.NonceSize()))
		if err != nil {
			return
		}
	} else {
		nonceParam, err = hex.DecodeString(nonce)
		if err != nil {
			return
		}
	}

	result = gcm.Seal(nil, nonceParam, []byte(originData), nil)
	return
}

// Encrypt AES用gcm模式解密
func (a AES) Decrypt(encryptedData string) (result []byte, err error) {
	return a.DecryptWithCustomData(encryptedData, DEFAULT_SECRETKEY, DEFAULT_SECRET_MESSAGE)
}

func (a AES) DecryptWithCustomData(encryptedData string, secretKey string, nonce string) (result []byte, err error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	ciphertext, err := hex.DecodeString(encryptedData)
	if err != nil {
		return
	}

	nonceParam, err := hex.DecodeString(nonce)
	if err != nil {
		return
	}

	result, err = gcm.Open(nil, nonceParam, ciphertext, nil)
	if err != nil {
		return
	}
	return
}

//AES key必须是16字节 24字节 32字节中的一种
func getValidAESKey(strKey string) []byte {
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
