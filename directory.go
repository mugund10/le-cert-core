package lecertcore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Stag string = "https://acme-staging-v02.api.letsencrypt.org/directory"
const Prod string = "https://acme-v02.api.letsencrypt.org/directory"

// le acme directory structure
type directory struct {
	Newnonce   string `json:"newNonce"`
	Newaccount string `json:"newAccount"`
	Neworder   string `json:"newOrder"`
	Revokecert string `json:"revokeCert"`
	Keychange  string `json:"keyChange"`
	Meta       meta   `json:"meta"`
	Renewal    string `json:"renewalInfo"`
}

type meta struct {
	Caaidentities []string `json:"caaIdentities"`
	Profiles      profiles `json:"profiles"`
	Tos           string   `json:"termsOfService"`
	Website       string   `json:"website"`
}

type profiles struct {
	Classic    string `json:"classic"`
	Shortlived string `json:"shortlived"`
	Tlsserver  string `json:"tlsserver"`
}

// makes gets requests to the acme server
// use "Stag" for staging url and
// "Prod" for production url).
func GetDir(url string) (*directory, error) {
	var Dir directory
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting response from %s", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data from response {%s}", resp.Body)
	}
	err = json.Unmarshal(body, &Dir)
	if err != nil {
		return nil, fmt.Errorf("error converting json | %s", err)
	}
	return &Dir, nil
}

// gets initial nonce from acme server
func (dir *directory) getNonce() (string, error) {
	resp, err := http.Head(dir.Newnonce)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	return resp.Header.Get("replay-nonce"), nil
}

