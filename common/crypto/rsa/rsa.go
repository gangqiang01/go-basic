package rsa

import (
	"bytes"
	"crypto/rand"
	crsa "crypto/rsa"
	"crypto/x509"

	"k8s.io/klog/v2"
)

func Encrypt(rawData []byte, publicKey []byte) ([]byte, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	pub := pubInterface.(*crsa.PublicKey)

	keySize, contentSize := pub.Size(), len(rawData)
	start := 0
	end := 0
	encryptLength := keySize - 11
	buffer := bytes.Buffer{}

	for start < contentSize {
		end = start + encryptLength
		if end > contentSize {
			end = contentSize
		}

		encryptBytes, err := crsa.EncryptPKCS1v15(rand.Reader, pub, rawData[start:end])
		if err != nil {
			klog.Errorf("err: %v", err)
			return nil, err
		}

		buffer.Write(encryptBytes)
		start = end
	}

	return buffer.Bytes(), nil
}

func Decrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	priv, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	priKey := priv.(*crsa.PrivateKey)

	keySize, textSize := priKey.Size(), len(ciphertext)
	start := 0
	end := 0
	buffer := bytes.Buffer{}

	for start < textSize {
		end = start + keySize
		if end > textSize {
			end = textSize
		}

		decryptBytes, err := crsa.DecryptPKCS1v15(rand.Reader, priKey, ciphertext[start:end])
		if err != nil {
			klog.Errorf("err: %v", err)
			return nil, err
		}

		buffer.Write(decryptBytes)
		start = end
	}

	return buffer.Bytes(), nil
}
