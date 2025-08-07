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
	return res.post(url, body, oResp)

}

// finalizes order
func (ordR *orderResp) Finalize(body io.Reader) (certbyte []byte, nonce string, location string, err error) {
	res := acme[orderResp]{res: *ordR}
	nonce, location, err = res.post(ordR.Final, body, ordR)
	if err != nil {
		return nil, "", "", err
	}
	if ok, err := res.poll(location, "order", ordR); ok {
		certbyte, err = res.get(ordR.Certificate)
		if err != nil {
			return nil, "", "", err
		}
		return certbyte, nonce, location, nil
	} else {
		return nil, "", "", err
	}
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
	err = json.Unmarshal(bod, authz)
	if err != nil {
		return &authResp{}, err
	}
	return authz, nil
}

func GetOrderResp() *orderResp {
	return &orderResp{}
}
