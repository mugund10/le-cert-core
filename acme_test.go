package lecertcore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	dom := []string{"homeserver.mugund10.top"}
	// keys generate
	kys, err := Loadkeys("accounttesting")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
	dir, err := GetDir(Stag)
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	non, err := dir.getNonce()
	if err != nil {
		t.Errorf("[nonce] %s", err)
	}
	// account create
	acc := NewAccount("bjmugundhan@gmail.com", true)
	finalout, err := JwsPayload(acc, *kys.private, non, dir.Newaccount, "")
	if err != nil {
		t.Errorf("[Jwspayload - account] %s", err)
	}
	non, kid, err := acc.Create(dir.Newaccount, strings.NewReader(string(finalout)))
	if err != nil {
		t.Errorf("[account] %s", err)
	}
	// submit order
	order := NewOrder(dom)
	finalout, err = JwsPayload(order, *kys.private, non, dir.Neworder, kid)
	if err != nil {
		t.Errorf("[Jwspayload - order] %s", err)
	}
	orresp := &orderResp{}
	non, orderloc, err := order.Submit(dir.Neworder, strings.NewReader(string(finalout)), orresp)
	if err != nil {
		t.Errorf("[order] %s", err)
	}
	//getting challenges
	authz, err := orresp.GetAuth()
	if err != nil {
		t.Errorf("[orderResp] %s", err)
	}
	non, _, err = authz.DohttpChallenge(*kys.private, non, kid)
	if err != nil {
		t.Errorf("[http chall] %s", err)
	}

	// sub csr
	kyss, err := Loadkeys("csr")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
	csr, err := kyss.GenCSR(dom)
	if err != nil {
		t.Errorf("[csr Gen] %s", err)
	}
	csrEnc := encodeToBase64(csr)
	build := csrDer{Csr: csrEnc}
	cout, err := JwsPayload(build, *kys.private, non, orresp.Final, kid)
	if err != nil {
		t.Errorf("[Jwspayload - csr] %s", err)
	}
	non1, orderloc1, err1 := orresp.finalize(strings.NewReader(string(cout)), orresp)
	if err1 != nil {
		t.Errorf("[finalize] %s", err)
	}
	log.Println(non1, orderloc1, err1)

	for {
		time.Sleep(2 * time.Second)
		resp2, err := http.Get(orderloc)
		if err != nil {
			t.Fatalf("[poll order] %s", err)
		}
		body2, err := io.ReadAll(resp2.Body)
		resp2.Body.Close()
		if err != nil {
			t.Fatalf("[read order poll] %s", err)
		}

		var updated orderResp
		if err := json.Unmarshal(body2, &updated); err != nil {
			t.Fatalf("[unmarshal order] %s", err)
		}

		log.Println("Order status:", updated.Status)
		if updated.Status == "valid" {
			certResp, err := http.Get(updated.Certificate)
			if err != nil {
				t.Fatalf("[download cert] %s", err)
			}
			certBody, err := io.ReadAll(certResp.Body)
			certResp.Body.Close()
			if err != nil {
				t.Fatalf("[read cert] %s", err)
			}
			err = SaveCertToFile("cert.pem", certBody)
			if err != nil {
				t.Fatalf("[save cert] %s", err)
			}
			log.Println("Certificate saved as cert.pem")

			// Save the private key
			err = SavePrivateKeyAsPEM("cert.key", kyss.private)
			if err != nil {
				t.Fatalf("[save key] %s", err)
			}
			log.Println("Private key saved as cert.key")
			break
		}

		if updated.Status == "invalid" {
			t.Fatalf("Order failed: %s", string(body2))
		}
	}

}

func SaveCertToFile(filename string, cert []byte) error {
	return os.WriteFile(filename, cert, 0644)
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
