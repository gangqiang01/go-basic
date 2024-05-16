package crypto

import (
	"encoding/base64"

	"github.com/edgehook/ithings/common/crypto/descbc"
	"github.com/edgehook/ithings/common/crypto/rsa"
	"k8s.io/klog/v2"
)

const (
	defaultKey = "ahc*5f/8"
)

/*
* Encrypt:
* 1. DES encrypt 2. Base64 encrypt.
 */
func Encrypt(data []byte) (string, error) {
	key := []byte(defaultKey)

	crypted, err := descbc.Encrypt(data, key)
	if err != nil {
		klog.Errorf("err: %v", err)
		return "", err
	}

	sData := base64.StdEncoding.EncodeToString(crypted)

	return sData, nil
}

/*
* Decrypt:
* 1. base64 decrypt 2. des decrypt.
 */
func Decrypt(s string) ([]byte, error) {
	key := []byte(defaultKey)

	crypted, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	data, err := descbc.Decrypt(crypted, key)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return data, nil
}

/*
* RSAEncrypt
* 1. RSA encrypt 2. Base64 encrypt.
 */
func RSAEncrypt(data []byte, cryptedPubKey string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(cryptedPubKey)
	if err != nil {
		klog.Errorf("err: %v", err)
		return "", err
	}

	encryptedBytes, err := rsa.Encrypt(data, publicKey)
	if err != nil {
		klog.Errorf("err: %v", err)
		return "", err
	}

	sData := base64.StdEncoding.EncodeToString(encryptedBytes)

	return sData, nil
}

/*
* RSADecrypt
* 1. RSA decrypt 2. Base64 decrypt.
 */
func RSADecrypt(data, cryptedPrivateKey string) ([]byte, error) {
	privateKey, err := base64.StdEncoding.DecodeString(cryptedPrivateKey)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return rsa.Decrypt(cipherText, privateKey)
}
