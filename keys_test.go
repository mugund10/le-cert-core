package lecertcore

import (
	"testing"
)

func TestCreateKeys(t *testing.T) {
	kys := CreateKeys()
	err := kys.Save("keytesting")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
}

func TestLoadkeys(t *testing.T) {
	_, err := Loadkeys("keytesting")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}

}
