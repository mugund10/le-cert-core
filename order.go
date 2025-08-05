package lecertcore

import (
	"encoding/json"
	"io"
	"net/http"
)

// ref : https://github.com/letsencrypt/boulder/blob/main/docs/acme-divergences.md
type order struct {
	Identifiers []identifier `json:"identifiers"`
}

type identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type orderResp struct {
	Status      string       `json:"status"`
	Expires     string       `json:"expires"`
	Identifiers []identifier `json:"identifiers"`
	Auth        []string     `json:"authorizations"`
	Final       string       `json:"finalize"`
	Certificate string       `json:"certificate,omitempty"`
}

func NewOrder(domains []string) order {
	ids := []identifier{}
	for _, domain := range domains {
		id := identifier{"dns", domain}
		ids = append(ids, id)
	}
	return order{ids}
}

// submits neworder
func (ord order) Submit(url string, body io.Reader, oResp *orderResp) (nonce string, location string, err error) {
	res := acme[order]{res: ord}
	nonce, location, err = res.post(url, body, oResp)
	return
}

func (ordR orderResp) finalize(body io.Reader, oResp *orderResp) (nonce string, location string, err error) {
	res := acme[orderResp]{res: ordR}
	nonce, location, err = res.post(ordR.Final, body, oResp)
	return
}

func (ordR orderResp) GetAuth() (*authResp, error) {
	resp, err := http.Get(ordR.Auth[0])
	if err != nil {
		return &authResp{}, err
	}
	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		return &authResp{}, err
	}
	authz := &authResp{}
	json.Unmarshal(bod, authz)
	return authz, nil
}
