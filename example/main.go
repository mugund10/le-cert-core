package main

import (
	"log"
	"strings"

	leacme "github.com/mugund10/le-cert-core"
)

func main() {
	dom := []string{"homeserver.mugund10.top"}
	// keys generate
	kname := "acc2"
	kys, err := leacme.Loadkeys(kname)
	if err != nil {
		log.Println("[keys] ", err)
		kys = leacme.CreateKeys()
		if err = kys.Save(kname); err != nil {
			log.Println("[keys] ", err)
		}
	}
	priv := kys.GetKeys()
	dir, err := leacme.GetDir(leacme.Stag)
	if err != nil {
		log.Println("[directory] ", err)
	}
	non, err := dir.GetNonce()
	if err != nil {
		log.Println("[nonce] ", err)
	}

	// account create
	acc := leacme.NewAccount("bjmugundhan@gmail.com", true)
	finalout, err := leacme.JwsPayload(acc, *priv, non, dir.Newaccount, "")
	if err != nil {
		log.Println("[Jwspayload - account] ", err)
	}
	non, kid, err := acc.Create(dir.Newaccount, strings.NewReader(string(finalout)))
	if err != nil {
		log.Println("[account] ", err)
	}
	log.Println("non", non, "kid", kid, "err", err)

	// submit order
	order := leacme.NewOrder(dom)
	finalout, err = leacme.JwsPayload(order, *priv, non, dir.Neworder, kid)
	if err != nil {
		log.Println("[Jwspayload - order] ", err)
	}
	orresp := leacme.GetOrderResp()
	non, orderloc, err := order.Submit(dir.Neworder, strings.NewReader(string(finalout)), orresp)
	if err != nil {
		log.Println("[order] ", err)
	}
	log.Println("non", non, "orderloc", orderloc, "err", err)

	//getting challenges
	authz, err := orresp.GetAuth()
	if err != nil {
		log.Println("[orderResp] ", err)
	}
	non, chalLoc, err := authz.DohttpChallenge(*priv, non, kid)
	if err != nil {
		log.Println("[http chall] ", err)
	}
	log.Println("non", non, "chalLoc", chalLoc, "err", err)

	// create csr
	build, err := leacme.CreateCsr("csr", dom)
	if err != nil {
		log.Println("[keys] ", err)
	}
	cout, err := leacme.JwsPayload(build, *priv, non, orresp.Final, kid)
	if err != nil {
		log.Println("[Jwspayload - csr] ", err)
	}
	cert, _, _, err1 := orresp.Finalize(strings.NewReader(string(cout)))
	if err1 != nil {
		log.Println("[finalize] ", err1)
	}
	//saving certs
	err = build.SaveCert(cert)
	if err1 != nil {
		log.Println(err)
	}

}
