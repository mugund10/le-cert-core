package lecertcore

import (
	"fmt"
	"testing"
)

func TestCreateKeys(t *testing.T) {
	kys := CreateKeys()
	fmt.Println(kys)
	err := kys.Save("testing")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
}

func TestLoadkeys(t *testing.T) {
	pkey, err := Loadkeys("testing")
	if err != nil {
		t.Errorf("[keys] %s", err)
	}
	fmt.Println(pkey, err)
}
