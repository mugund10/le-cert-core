package lecertcore

import (
	"testing"
)

func TestCreateKeys(t *testing.T) {
	kys := CreateKeys()
	err := kys.Save("accounttesting")
	if err != nil {
		t.Errorf("[keys / create] %s", err)
	}
}

func TestLoadkeys(t *testing.T) {
	_, err := Loadkeys("accounttesting")
	if err != nil {
		t.Errorf("[keys / load] %s", err)
	}

}

func TestCreateCsr(t *testing.T) {
	dom := []string{"domain"}
	_, err := CreateCsr("accounttesting", dom)
	if err != nil {
		t.Errorf("[keys / createcsr] %s", err)
	}

}
