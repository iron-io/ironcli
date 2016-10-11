package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func RsaEncrypt(publicKey []byte, payloadPlain string) (string, error) {
	rsablock, _ := pem.Decode(publicKey)
	if rsablock == nil {
		return "", errors.New("failed to find any PEM data in key input")
	}

	rsaKey, err := x509.ParsePKIXPublicKey(rsablock.Bytes)
	if err != nil {
		return "", err
	}
	rsaPublicKey := rsaKey.(*rsa.PublicKey)

	// get a random aes-128 session key to encrypt
	aesKey := make([]byte, 128/8)
	if _, err := rand.Read(aesKey); err != nil {
		return "", err
	}

	// have to use sha1 b/c ruby openssl picks it for OAEP:  https://www.openssl.org/docs/manmaster/crypto/RSA_public_encrypt.html
	aesKeyCipher, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, rsaPublicKey, aesKey, nil)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	pbytes := []byte(payloadPlain)
	// The IV needs to be unique, but not secure. last 12 bytes are IV.
	ciphertext := make([]byte, len(pbytes)+gcm.Overhead()+gcm.NonceSize())
	nonce := ciphertext[len(ciphertext)-gcm.NonceSize():]
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	// tag is appended to cipher as last 16 bytes. https://golang.org/src/crypto/cipher/gcm.go?s=2318:2357#L145
	gcm.Seal(ciphertext[:0], nonce, pbytes, nil)
	return base64.StdEncoding.EncodeToString(append(aesKeyCipher, ciphertext...)), nil
}
