package lecertcore

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type authResp struct {
	Identifier []identifier `json:"identifier"`
	Status     string       `json:"status"`
	Expires    string       `json:"expires"`
	Challenges []challenge  `json:"challenges"`
}

type challenge struct {
	Typ    string `json:"type"`
	Url    string `json:"url"`
	Status string `json:"pending"`
	Token  string `json:"token"`
}

type challengeResp struct {
	Status string `json:"status"`
}

func (authz *authResp) DohttpChallenge(pkey ecdsa.PrivateKey, nonce string, kid string) (non string, chalLoc string, err error) {
	for _, chall := range authz.Challenges {
		if chall.Typ == "http-01" {
			byt, err := JwsPayload(nil, pkey, nonce, chall.Url, kid)
			if err != nil {
				return "", "", err
			}
			thumb := ThumbPrint(pkey.PublicKey)
			go func(chall challenge, thumb string) {
				http.HandleFunc("/.well-known/acme-challenge/"+chall.Token, func(w http.ResponseWriter, r *http.Request) {
					keyAuth := chall.Token + "." + thumb
					fmt.Fprint(w, keyAuth)
				})
				log.Fatal(http.ListenAndServe(":80", nil))
			}(chall, thumb)
			time.Sleep(time.Second * 6)
			non, chalLoc, err := chall.respond(chall.Url, strings.NewReader(string(byt)))
			if err != nil {
				return "", "", err
			} else {
				return non, chalLoc, err
			}

		}
	}
	return "", "", fmt.Errorf("no http challenge found")
}

// submits neworder
func (chall challenge) respond(url string, body io.Reader) (nonce string, location string, err error) {
	res := acme[challenge]{res: chall}
	return res.post(url, body, &challengeResp{})
}
