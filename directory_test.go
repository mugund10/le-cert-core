package lecertcore

import (
	"testing"
)

func TestGetNonce(t *testing.T) {
	dir, err := GetDir()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
	_, err = dir.GetNonce()
	if err != nil {
		t.Errorf("[directory] %s", err)
	}
}
