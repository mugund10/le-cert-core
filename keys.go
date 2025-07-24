package lecertcore

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type keys struct {
	public  *ed25519.PublicKey
	private *ed25519.PrivateKey
}

// Creates Private and Public keys
func CreateKeys() keys {
	pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return keys{&pubkey, &privkey}
}

// Saves Keys to the file
func (k *keys) Save(filename string) error {
	if k.prConvAndSave(filename) && k.puConvAndSave(filename) {
		return nil
	} else {
		return fmt.Errorf("[error] Saving Keys")
	}
}

// marshals privates key
func (k *keys) prConvAndSave(filename string) bool {
	privbyte, err := x509.MarshalPKCS8PrivateKey(*k.private)
	if err != nil {
		log.Println(err)
		return false
	}
	pBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privbyte,
	}
	pemData := pem.EncodeToMemory(pBlock)
	err = saveAsFile(fmt.Sprintf("%s_privatekey.pem", filename), pemData)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// marshals pubic key
func (k *keys) puConvAndSave(filename string) bool {
	publicbyte, err := x509.MarshalPKIXPublicKey(*k.public)
	if err != nil {
		log.Println(err)
		return false
	}
	pBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicbyte,
	}
	pemData := pem.EncodeToMemory(pBlock)
	err = saveAsFile(fmt.Sprintf("%s_publickey.pem", filename), pemData)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// saves file with the given name
func saveAsFile(name string, byt []byte) error {
	err := os.WriteFile(name, byt, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Loads keys from pem files
func Loadkeys(filename string) (keys, error) {
	privPEM, err := readfile(fmt.Sprintf("%s_privatekey.pem", filename))
	if err != nil {
		return keys{}, err
	}
	prkey, err := x509.ParsePKCS8PrivateKey(privPEM.Bytes)
	if err != nil {
		return keys{}, err
	}
	priv, ok := prkey.(ed25519.PrivateKey)
	if !ok {
		return keys{}, fmt.Errorf("[error] Its not an ed25519 privatekey")
	}

	pubPEM, err := readfile(fmt.Sprintf("%s_publickey.pem", filename))
	if err != nil {
		return keys{}, err
	}
	pukey, err := x509.ParsePKIXPublicKey(pubPEM.Bytes)
	if err != nil {
		return keys{}, err
	}
	pub, ok := pukey.(ed25519.PublicKey)
	if !ok {
		return keys{}, fmt.Errorf("[error] Its not an ed25519 publickey")
	}

	return keys{&pub, &priv}, nil
}

// reads pem files and return decode pem block
func readfile(filename string) (*pem.Block, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("[error] reading pem file : %s", filename)
	}
	pblock, _ := pem.Decode(byt)
	if pblock == nil {
		return nil, fmt.Errorf("[error] failed to decode PEM block from file : %s", filename)
	}
	return pblock, nil
}

func (k *keys) GetKeys() (ed25519.PublicKey, ed25519.PrivateKey) {
	return *k.public, *k.private
}
