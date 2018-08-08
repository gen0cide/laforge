package buildutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// GenerateRSAPrivateKey creates a RSA Private Key of specified byte size
func GenerateRSAPrivateKey(size int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// EncodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}
	return pem.EncodeToMemory(&privBlock)
}

// EncodePublicKeyToSSH takes a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func EncodePublicKeyToSSH(pubkey *rsa.PublicKey) ([]byte, error) {
	publicRSAKey, err := ssh.NewPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRSAKey)
	return pubKeyBytes, nil
}

// GenerateSSHKeyPair is a small helper function that just returns an SSH key pair of a given size easily
func GenerateSSHKeyPair(size int) (privkey string, pubkey string, err error) {
	privk, err := GenerateRSAPrivateKey(size)
	if err != nil {
		return
	}
	privkd := EncodePrivateKeyToPEM(privk)
	pubkd, err := EncodePublicKeyToSSH(&privk.PublicKey)
	if err != nil {
		return
	}
	privkey = string(privkd)
	pubkey = string(pubkd)
	return
}

// WriteKeyfile writes keys to a file
func WriteKeyfile(keyBytes []byte, saveFileTo string) error {
	return ioutil.WriteFile(saveFileTo, keyBytes, 0600)
}
