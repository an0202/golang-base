/**
 * @Author: jie.an
 * @Description:
 * @File:  hash.go
 * @Version: 1.0.0
 * @Date: 2020/3/22 1:23 下午
 */

package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func SHA256(data string) (HashedString string) {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// https://www.jianshu.com/p/43820e5c08c3
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

func EncryptAES(srcString, keyString string) (string, error) {
	src := []byte(srcString)
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error(), err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return string(src[:]), nil
}

func DecryptAES(srcString, keyString string) (string, error) {
	src := []byte(srcString)
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error(), err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return string(src[:]), nil
}

func DecodeBase64String(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}
