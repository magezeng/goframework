package Crypto

import (
	"encoding/hex"
	"testing"
)

var (
	aesInstance = AES{}
)

func TestAES_EncryptAndDecrypt(t *testing.T) {
	encryptedData, err := aesInstance.Encrypt("node1 45.22.11.232:11919 2KvQbBm2QRLpUwoyR9oB6nq4yRV78AmC72wMh9nSsFGb1QANLKN b140bb1f9c2b098gf269c354c740d727afbc9dd4254afc23f52c1e4b828a4df2 2")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("[加密后的结果] %x", encryptedData)
	encryptedStr := hex.EncodeToString(encryptedData)
	hexByte, err := hex.DecodeString(encryptedStr)
	if err != nil {
		t.Error(err)
		t.Fatal()
	}
	result, err := aesInstance.Decrypt(string(hexByte))
	if err != nil {
		t.Error(err)
		t.Fatal()
	}
	t.Logf("[解密后的字符串] %s", string(result))
}

func TestAES_DecryptWithCustomData(t *testing.T) {
	encryptedData := "8656e959df0121420801abccb779e41dc4720fcf751493c5cbc9f86499ca9edddac205effe356ea79284f474c2e6bfa5444285475d5ac34749d42b2158f185"
	result, err := aesInstance.DecryptWithCustomData(encryptedData, "AES256Key-32Characters1234567890", "bb8ef84243d2ee95a41c6c57")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(result))
}
