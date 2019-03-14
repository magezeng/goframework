package Crypto

import "testing"

var (
	aesInstance = AES{}
)

func TestAES_EncryptAndDecrypt(t *testing.T) {
	encryptedData, err := aesInstance.Encrypt("123456")
	if err != nil{
		t.Error(err)
		t.Fail()
	}
	t.Logf("[加密后的结果] %x", encryptedData)
	result, err := aesInstance.Decrypt(string(encryptedData))
	if err != nil{
		t.Error(err)
		t.Fail()
	}
	t.Logf("[解密后的字符串] %s", string(result))
}

