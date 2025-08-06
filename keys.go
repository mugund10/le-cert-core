package lecertcore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type keys struct {
	private *ecdsa.PrivateKey
}

type csrDer struct {
	Csr string `json:"csr"`
}

// Creates Private and Public keys
func CreateKeys() keys {
	privkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return keys{privkey}
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

	privbyte, err := x509.MarshalPKCS8PrivateKey(k.private)
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
	publicbyte, err := x509.MarshalPKIXPublicKey(&k.private.PublicKey)
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
	return os.WriteFile(name, byt, 0644)
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
	priv, ok := prkey.(*ecdsa.PrivateKey)
	if !ok {
		return keys{}, fmt.Errorf("[error] Its not an ecdsa privatekey")
	}

	// pubPEM, err := readfile(fmt.Sprintf("%s_publickey.pem", filename))
	// if err != nil {
	// 	return keys{}, err
	// }
	// pukey, err := x509.ParsePKIXPublicKey(pubPEM.Bytes)
	// if err != nil {
	// 	return keys{}, err
	// }
	// pub, ok := pukey.(ed25519.PublicKey)
	// if !ok {
	// 	return keys{}, fmt.Errorf("[error] Its not an ed25519 publickey")
	// }

	return keys{priv}, nil
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

func (k *keys) GetKeys() *ecdsa.PrivateKey {
	return k.private
}

func (k *keys) GenCSR(domains []string) ([]byte, error) {
	cont := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domains[0],
		},
		DNSNames: domains,
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, &cont, k.private)
	if err != nil {
		return nil, err
	}
	return der, nil
}

func SavePrivateKeyAsPEM(filename string, key *ecdsa.PrivateKey) error {
	privBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, block)
}
