package lecertcore

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestNewAccount(t *testing.T) {
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
		t.Errorf("[directory] %s", err)
	}
	acc := NewAccount("bjmugundhan@gmail.com", true)
	ajws := NewJws(acc, kys.private.PublicKey, non, dir.Newaccount)
	flat := ajws.EncodeFlat(kys.private)
	log.Println("acc", flat)
	finalout, err := json.Marshal(flat)
	log.Println("acc", string(finalout))
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	acc.Register(dir.Newaccount, strings.NewReader(string(finalout)))

}
