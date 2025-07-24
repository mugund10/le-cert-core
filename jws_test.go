package lecertcore

import (
	"log"
	"testing"
)

func TestNewJws(t *testing.T) {
	keys, err := Loadkeys("testing")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
	payload := map[string]interface{}{
		"sub":   "1234567890",
		"name":  "John Doe",
		"admin": true,
		"iat":   1516239022,
	}

	ajws := NewJws(payload)
	out := ajws.Gen(keys.private)
	log.Println(string(out))
}
