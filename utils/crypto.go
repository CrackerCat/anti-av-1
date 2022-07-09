package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func Crypt(key, src []byte) ([]byte, error) {
	k := hex.EncodeToString(key[:8])
	x, err := aes.NewCipher([]byte(k))
	if err != nil {
		return nil, err
	}
	src = PKCS7Padding(src, x.BlockSize())
	kl := len(key)
	dst := make([]byte, len(src)+kl)
	copy(dst, key)
	mode := cipher.NewCBCEncrypter(x, []byte(k))
	mode.CryptBlocks(dst[kl:], src)
	return dst, nil
}

func DeCrypt(src []byte) ([]byte, error) {
	key := hex.EncodeToString(src[:8])
	src = src[8:]
	x, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(x, []byte(key))
	dst := make([]byte, len(src))
	mode.CryptBlocks(dst, src)
	dst = PKCS7UnPadding(dst)
	return dst, nil
}

//PKCS7Padding say ...
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

//PKCS7UnPadding 使用PKCS7进行填充 复原
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	unPadding := int(origData[length-1])
	if unPadding > length {
		return origData
	}
	return origData[:(length - unPadding)]
}
