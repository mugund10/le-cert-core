package lecertcore

import (
	"fmt"
	"io"
)

const contTyp = "application/jose+json"

type account struct {
	Status  string   `json:"status"`
	Contact []string `json:"contact"`
	Tos     bool     `json:"termsOfServiceAgreed"`
	//exists  bool   `json:"onlyReturnExisting"`
}

type accountResp struct {
}

// creates a new account
func NewAccount(mailaddr string, tos bool) *account {
	return &account{
		Status:  "valid",
		Contact: []string{fmt.Sprintf("mailto:%s", mailaddr)},
		Tos:     tos,
	}
}

// creates new account
func (a account) Create(url string, body io.Reader) (nonce string, location string, err error) {
	res := acme[account]{res: a}
	return res.post(url, body, &accountResp{})
}
