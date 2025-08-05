package lecertcore

import (
	"encoding/json"
	"io"
	"net/http"
)

type acme[resource account | order | challenge | orderResp] struct {
	res resource
}

func (r acme[res]) post(url string, body io.Reader, acresp any) (nonce string, location string, err error) {
	resp, err := http.Post(url, contTyp, body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		nonce = resp.Header.Get("Replay-Nonce")
		location = resp.Header.Get("Location")
		if err = json.Unmarshal(bod, acresp); err != nil {
			return
		}
		return
	}
	return
}

func (r acme[res]) get() {

}
