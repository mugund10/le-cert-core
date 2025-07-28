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

// signs head and payload with ed25519 privatekey
func sign(proc, pl string, privKey *ecdsa.PrivateKey) string {
	msg := []byte(proc + "." + pl)
	hash := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		log.Println(err)
	}
	sig := append(EncodeCoOrd(r), EncodeCoOrd(s)...)
	return encodeToBase64(sig)
}

// encodes content to base64
func encodeToBase64(content []byte) string {
	return base64.RawURLEncoding.EncodeToString(content)
}

// ref: https://datatracker.ietf.org/doc/html/rfc8555/#section-6.2
type protected struct {
	Alg   string    `json:"alg"`
	Typ   string    `json:"typ,omitempty"`
	Nonce string    `json:"nonce"`
	Url   string    `json:"url"`
	Jwk   jwkHeader `json:"jwk,omitempty"`
	Kid   string    `json:"kid,omitempty"`
}

// ref: https://datatracker.ietf.org/doc/html/rfc7518#section-6.2
type jwkHeader struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	D   string `json:"d,omitempty"` // for private key
	X   string `json:"x"`
	Y   string `json:"y"`
}

type payload struct {
	any
}

// adds basic structure to the jws
func NewJws(pload any, pbkey ecdsa.PublicKey, nonce string, url string) *jws {
	//pubkey := encodeToBase64(pbkey)
	jwk := jwkHeader{Kty: "EC", Crv: "P-256", X: encodeToBase64(EncodeCoOrd(pbkey.X)), Y: encodeToBase64(EncodeCoOrd(pbkey.Y))}
	//ref: https://datatracker.ietf.org/doc/html/rfc7515#appendix-A.3
	proc := protected{Alg: "ES256", Typ: "JWS", Url: url, Nonce: nonce, Jwk: jwk}
	pl := payload{pload}
	return &jws{payload: pl, protected: proc}
}

// supports only p256
func EncodeCoOrd(coOrd *big.Int) []byte {
	coBytes := coOrd.Bytes()
	coPad := make([]byte, 32)
	copy(coPad[32-len(coBytes):], coBytes)
	return coPad
}

func (j *jws) addreplayNonce(nonce string) {
	j.protected.Nonce = nonce
}
