package vemcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
)

// GenerateECDHKeypair generates a P-256 ECDH keypair
func GenerateECDHKeypair() (*ecdh.PrivateKey, string, error) {
	curve := ecdh.P256()

	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(priv.PublicKey())
	if err != nil {
		return nil, "", err
	}

	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)
	return priv, pubBase64, nil
}

type EncryptedVemResult struct {
	EphemeralPubKey string `json:"ephemeralPubKey"`
	Nonce           string `json:"nonce"`
	Ciphertext      string `json:"ciphertext"`
}

// DecryptVemResult decrypts the ECIES payload
func DecryptVemResult(
	priv *ecdh.PrivateKey,
	encryptedBase64 string,
) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, err
	}

	var enc EncryptedVemResult
	if err := json.Unmarshal(raw, &enc); err != nil {
		return nil, err
	}

	ephemeralPubBytes, _ := base64.StdEncoding.DecodeString(enc.EphemeralPubKey)
	nonce, _ := base64.StdEncoding.DecodeString(enc.Nonce)
	ciphertext, _ := base64.StdEncoding.DecodeString(enc.Ciphertext)

	ephemeralPub, err := ecdh.P256().NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := priv.ECDH(ephemeralPub)
	if err != nil {
		return nil, err
	}

	key := sha256.Sum256(sharedSecret)

	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)

	return gcm.Open(nil, nonce, ciphertext, nil)
}
