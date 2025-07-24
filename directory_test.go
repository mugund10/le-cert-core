package lecertcore

import (
	"testing"
)

func TestGetNonce(t *testing.T) {
	dir, err := GetDir(Stag)
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	_, err = dir.getNonce()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
}
