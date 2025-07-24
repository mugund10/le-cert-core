package lecertcore

type Account struct {
	Contact string `json:"contact"`
	Tos     bool   `json:"termsOfServiceAgreed"`
	//exists  bool   `json:"onlyReturnExisting"`
}
