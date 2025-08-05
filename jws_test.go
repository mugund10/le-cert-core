package lecertcore

import (
	"log"
	"testing"
)

func TestNewJws(t *testing.T) {
	keys, err := Loadkeys("jwstesting")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
	payload := map[string]interface{}{
		"sub":   "1234567890",
		"name":  "John Doe",
		"admin": true,
		"iat":   1516239022,
	}
	dir, err := GetDir(Stag)
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	non, err := dir.getNonce()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	log.Println("nonce:", non)
	ajws := newJws(payload, keys.private.PublicKey, non, "randomurl", "")
	out := ajws.EncodeFlat(keys.private)
	log.Println(out)
}
