package Crypto

type CryptoInterface interface {
	Encrypt(originData string) (result []byte, err error)
	Decrypt(encryptedData string) (result []byte, err error)
}
