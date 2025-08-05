package lecertcore

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
)

// ref: https://datatracker.ietf.org/doc/html/rfc7515#appendix-A.7
// unprotected header not allowed in acme. ref:https://datatracker.ietf.org/doc/html/rfc8555/#section-6.2
type flattenedJws struct {
	Payload   string `json:"payload"`
	Protected string `json:"protected"`
	Signature string `json:"signature"`
}

type jws struct {
	payload
	protected
}

// encode jws using ed25519 privatekey
func (j *jws) EncodeFlat(privKey *ecdsa.PrivateKey) flattenedJws {
	proc, err := json.Marshal(j.protected)
	if err != nil {
		log.Println(err)
	}
	pl, err := json.Marshal(j.payload.any)
	if err != nil {
		log.Println(err)
	}
	EncProc := encodeToBase64(proc)
	EncPl := encodeToBase64(pl)
	signature := sign(EncProc, EncPl, privKey)
	// jwsString := EncHead + "." + EncPl + "." + signature
	return flattenedJws{Payload: EncPl, Protected: EncProc, Signature: signature}
}

// Generates Thumbprint of Jwk
func ThumbPrint(pbkey ecdsa.PublicKey) string {
	tp := map[string]string{
		"crv": "P-256",
		"kty": "EC",
		"x":   encodeToBase64(encCoOrd(pbkey.X)),
		"y":   encodeToBase64(encCoOrd(pbkey.Y)),
	}
	marshaled, err := json.Marshal(tp)
	if err != nil {
		log.Printf("ThumbPrint: failed to marshal JWK: %v", err)
		return ""
	}
	hash := sha256.Sum256(marshaled)
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// signs head and payload with ed25519 privatekey
func sign(proc, pl string, privKey *ecdsa.PrivateKey) string {
	msg := []byte(proc + "." + pl)
	hash := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		log.Println(err)
	}
	sig := append(encCoOrd(r), encCoOrd(s)...)
	return encodeToBase64(sig)
}

// encodes content to base64
func encodeToBase64(content []byte) string {
	return base64.RawURLEncoding.EncodeToString(content)
}

// ref: https://datatracker.ietf.org/doc/html/rfc8555/#section-6.2
type protected struct {
	Alg   string     `json:"alg"`
	Typ   string     `json:"typ,omitempty"`
	Nonce string     `json:"nonce"`
	Url   string     `json:"url"`
	Jwk   *jwkHeader `json:"jwk,omitempty"`
	Kid   string     `json:"kid,omitempty"`
}

// ref: https://datatracker.ietf.org/doc/html/rfc7518#section-6.2
type jwkHeader struct {
	Crv string `json:"crv"`
	D   string `json:"d,omitempty"` // for private key
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type payload struct {
	any
}

// adds basic structure to the jws
// kid should be empty "" for the first request
func newJws(pload any, pbkey ecdsa.PublicKey, nonce string, url string, kid string) *jws {
	var jwk *jwkHeader
	//pubkey := encodeToBase64(pbkey)
	if kid == "" {
		jwk = &jwkHeader{Kty: "EC", Crv: "P-256", X: encodeToBase64(encCoOrd(pbkey.X)), Y: encodeToBase64(encCoOrd(pbkey.Y))}
	}
	//ref: https://datatracker.ietf.org/doc/html/rfc7515#appendix-A.3
	proc := protected{Alg: "ES256", Typ: "JWS", Url: url, Nonce: nonce, Jwk: jwk, Kid: kid}
	pl := payload{pload}
	return &jws{payload: pl, protected: proc}
}

// supports only p256
func encCoOrd(coOrd *big.Int) []byte {
	coBytes := coOrd.Bytes()
	coPad := make([]byte, 32)
	copy(coPad[32-len(coBytes):], coBytes)
	return coPad
}

func JwsPayload(pload any, pbkey ecdsa.PrivateKey, nonce string, url string, kid string) ([]byte, error) {
	ajws := newJws(pload, pbkey.PublicKey, nonce, url, kid)
	flat := ajws.EncodeFlat(&pbkey)
	byt, err := json.Marshal(flat)
	if err != nil {
		return nil, err
	}
	return byt, nil
}
