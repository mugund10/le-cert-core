package lecertcore

import "fmt"

type Account struct {
	Status  string   `json:"status"`
	Contact []string `json:"contact"`
	Tos     bool     `json:"termsOfServiceAgreed"`
	//exists  bool   `json:"onlyReturnExisting"`
}

// creates a new account
func NewAccount(mailaddr string, tos bool) *Account {
	return &Account{
		Status: "valid", 
		Contact: []string{fmt.Sprintf("mailto:%s", mailaddr)},
		Tos: tos,
	}


}
