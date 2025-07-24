package lecertcore

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"log"
)

type jws struct {
	header
	payload
	sign string
}

func (j *jws) Gen(privKey ed25519.PrivateKey) string {
	header, err := j.header.marshal()
	if err != nil {
		log.Println(err)
	}
	pl, err := j.payload.marshal()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(pl))
	EncHead := Encode(header)
	EncPl := Encode(pl)
	signature := sign(EncHead, EncPl, privKey)
	jwsString := EncHead + "." + EncPl + "." + signature
	return jwsString
}

func sign(head, pl string, privKey ed25519.PrivateKey) string {
	msg := []byte(head + "." + pl)
	sig := ed25519.Sign(privKey, msg)
	return Encode(sig)
}

func Encode(content []byte) string {
	return base64.RawURLEncoding.EncodeToString(content)
}

type header struct {
	Algorith string `json:"alg"`
	Typ      string `json:"typ"`
}

func (h header) marshal() ([]byte, error) {
	return json.Marshal(h)
}

type payload struct {
	any
}

func (p payload) marshal() ([]byte, error) {
	return json.Marshal(p.any)
}

// adds basic structure to the jws
func NewJws(pload any) *jws {
	head := header{Algorith: "EdDSA", Typ: "JWT"}
	pl := payload{pload}
	return &jws{head, pl, ""}
}
