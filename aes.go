package common_go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// PKCS7Padding 将明文填充为块长度的整数倍
func PKCS7Padding(p []byte, blockSize int) []byte {
	pad := blockSize - len(p)%blockSize
	padtext := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(p, padtext...)
}

// PKCS7UnPadding 从明文尾部删除填充数据
func PKCS7UnPadding(p []byte) []byte {
	length := len(p)
	paddLen := int(p[length-1])
	return p[:(length - paddLen)]
}

// AESCBCEncrypt CBC模式下用AES算法加密数据
// 请注意，密钥长度必须为16、24或32字节才能选择AES-128、AES-192或AES-256
// 请注意，AES块大小为16字节
func AESCBCEncrypt(p, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	p = PKCS7Padding(p, block.BlockSize())
	ciphertext := make([]byte, len(p))
	if len(iv) == 0 {
		iv = key[:block.BlockSize()]
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(ciphertext, p)
	return ciphertext, nil
}

// AESCBCDecrypt CBC模式下AES算法对密文的解密
// 请注意，密钥长度必须为16、24或32字节才能选择AES-128、AES-192或AES-256
// 请注意，AES块大小为16字节
func AESCBCDecrypt(c, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(c))
	if len(iv) == 0 {
		iv = key[:block.BlockSize()]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plaintext, c)
	return PKCS7UnPadding(plaintext), nil
}

// Base64AESCBCEncrypt 在CBC模式下使用AES算法对数据进行加密，并使用base64进行编码
// php对应代码 base64_encode(openssl_encrypt($p, 'AES-256-CBC', $key, OPENSSL_RAW_DATA, $iv));
// 请注意，密钥长度必须为16、24或32字节才能选择AES-128、AES-192或AES-256
// 请注意，AES块大小为16字节,iv必须为16位
func Base64AESCBCEncrypt(p, key, iv string) (string, error) {
	c, err := AESCBCEncrypt([]byte(p), []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(c), nil
}

// Base64AESCBCDecrypt CBC模式下用AES算法解密base64编码的密文
// php对应代码 openssl_decrypt(base64_decode($p), 'AES-256-CBC', $key, OPENSSL_RAW_DATA, $iv);
// 请注意，密钥长度必须为16、24或32字节才能选择AES-128、AES-192或AES-256
// 请注意，AES块大小为16字节,iv必须为16位
func Base64AESCBCDecrypt(c, key, iv string) ([]byte, error) {
	oriCipher, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return nil, err
	}
	p, err := AESCBCDecrypt(oriCipher, []byte(key), []byte(iv))
	if err != nil {
		return nil, err
	}
	return p, nil
}

// 示例代码
//	p := []byte("plaintext")
//	key := []byte("12345678abcdefgh")
//
//	ciphertext, _ := Base64AESCBCEncrypt(p, key)
//	fmt.Println(ciphertext)
//
//	plaintext, _ := Base64AESCBCDecrypt(ciphertext, key)
//	fmt.Println(string(plaintext))
//输出:
//A67NhD3RBiNaMgG6HTm8LQ==
//plaintext
