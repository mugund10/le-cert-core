package lecertcore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

type keys struct {
	private *ecdsa.PrivateKey
	public  *ecdsa.PublicKey
}

// Creates Private and Public keys
func CreateKeys() keys {
	pkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return keys{pkey, &pkey.PublicKey}
}

// Saves Keys to the file
func (k *keys) Save(filename string) error {
	if k.prConvAndSave(filename) && k.puConvAndSave(filename) {
		return nil
	} else {
		return fmt.Errorf("[error] Saving Keys")
	}
}

func (k *keys) prConvAndSave(filename string) bool {
	privbyte, err := x509.MarshalPKCS8PrivateKey(k.private)
	if err != nil {
		log.Println(err)
		return false
	}
	err = saveAsFile(fmt.Sprintf("%s_privatekey.pem", filename), privbyte)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (k *keys) puConvAndSave(filename string) bool {
	publicbyte, err := x509.MarshalPKIXPublicKey(k.public)
	if err != nil {
		log.Println(err)
		return false
	}
	err = saveAsFile(fmt.Sprintf("%s_publickey.pem", filename), publicbyte)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func saveAsFile(name string, byt []byte) error {
	err := os.WriteFile(name, byt, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Loads Keys from file
func Loadkeys() {
	x509.ParsePKCS8PrivateKey()
	x509.ParsePKIXPublicKey()
}
