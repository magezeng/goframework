package Crypto

import (
	"encoding/hex"
	"testing"
)

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
	encryptedStr := hex.EncodeToString(encryptedData)
	hexByte, err := hex.DecodeString(encryptedStr)
	if err != nil{
		t.Error(err)
		t.Fatal()
	}
	result, err := aesInstance.Decrypt(string(hexByte))
	if err != nil{
		t.Error(err)
		t.Fatal()
	}
	t.Logf("[解密后的字符串] %s", string(result))
}
