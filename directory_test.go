package lecertcore

import (
	"fmt"
	"testing"
)

func TestGetNonce(t *testing.T) {
	dir, err := GetDir()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	out, err := dir.GetNonce()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	fmt.Println("out = ", out)
}
