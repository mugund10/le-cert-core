package lecertcore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type acme[resource account | order | challenge | orderResp] struct {
	res resource
}

func (r acme[res]) post(url string, body io.Reader, acresp any) (nonce string, location string, err error) {
	resp, err := http.Post(url, contTyp, body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
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

func (r acme[res]) get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (r acme[res]) poll(url string, typ string, ores any) (bool, error) {
	switch typ {
	case "order":
		{
			ordresp, ok := ores.(*orderResp)
			if !ok {
				return false, fmt.Errorf("it should be orderresp")
			}
			for {
				time.Sleep(2 * time.Second)
				resp, err := r.get(url)
				if err != nil {
					return false, err
				}
				if err := json.Unmarshal(resp, ordresp); err != nil {
					return false, err
				}
				switch ordresp.Status {
				case "valid":
					return true, nil
				case "invalid":
					return false, nil
				}
			}
		}
	default:
		return false, fmt.Errorf("use order or challenge in the field typ")
	}

}
