package lecertcore

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const contTyp = "application/jose+json"

type Account struct {
	Status  string   `json:"status"`
	Contact []string `json:"contact"`
	Tos     bool     `json:"termsOfServiceAgreed"`
	//exists  bool   `json:"onlyReturnExisting"`
}

// creates a new account
func NewAccount(mailaddr string, tos bool) *Account {
	return &Account{
		Status:  "valid",
		Contact: []string{fmt.Sprintf("mailto:%s", mailaddr)},
		Tos:     tos,
	}
}

func (a *Account) Register(url string, body io.Reader) {
	resp, err := http.Post(url, contTyp, body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.StatusCode)
	log.Println("account body :", string(bod))
	log.Println(resp.Header.Get("Replay-Nonce"))
	log.Println(resp.Header.Get("Location"))
}
