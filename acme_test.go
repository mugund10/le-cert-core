package lecertcore

// func TestNewAccount(t *testing.T) {
// 	dom := []string{"homeserver.mugund10.top"}
// 	// keys generate
// 	kys, err := Loadkeys("acc2")
// 	if err != nil {
// 		t.Errorf("[keys] %s", err)
// 	}
// 	dir, err := GetDir(Stag)
// 	if err != nil {
// 		t.Errorf("[directory] %s", err)
// 	}
// 	non, err := dir.getNonce()
// 	if err != nil {
// 		t.Errorf("[nonce] %s", err)
// 	}

// 	// account create
// 	acc := NewAccount("bjmugundhan@gmail.com", true)
// 	finalout, err := JwsPayload(acc, *kys.private, non, dir.Newaccount, "")
// 	if err != nil {
// 		t.Errorf("[Jwspayload - account] %s", err)
// 	}
// 	non, kid, err := acc.Create(dir.Newaccount, strings.NewReader(string(finalout)))
// 	if err != nil {
// 		t.Errorf("[account] %s", err)
// 	}
// 	log.Println("non", non, "kid", kid, "err", err)

// 	// submit order
// 	order := NewOrder(dom)
// 	finalout, err = JwsPayload(order, *kys.private, non, dir.Neworder, kid)
// 	if err != nil {
// 		t.Errorf("[Jwspayload - order] %s", err)
// 	}
// 	orresp := &orderResp{}
// 	non, orderloc, err := order.Submit(dir.Neworder, strings.NewReader(string(finalout)), orresp)
// 	if err != nil {
// 		t.Errorf("[order] %s", err)
// 	}
// 	log.Println("non", non, "orderloc", orderloc, "err", err)

// 	//getting challenges
// 	authz, err := orresp.GetAuth()
// 	if err != nil {
// 		t.Errorf("[orderResp] %s", err)
// 	}
// 	non, chalLoc, err := authz.DohttpChallenge(*kys.private, non, kid)
// 	if err != nil {
// 		t.Errorf("[http chall] %s", err)
// 	}
// 	log.Println("non", non, "chalLoc", chalLoc, "err", err)

// 	// sub csr
// 	kyss, err := Loadkeys("csr")
// 	if err != nil {
// 		t.Errorf("[keys] %s", err)
// 	}
// 	csr, err := kyss.GenCSR(dom)
// 	if err != nil {
// 		t.Errorf("[csr Gen] %s", err)
// 	}
// 	csrEnc := encodeToBase64(csr)
// 	build := csrDer{Csr: csrEnc}
// 	cout, err := JwsPayload(build, *kys.private, non, orresp.Final, kid)
// 	if err != nil {
// 		t.Errorf("[Jwspayload - csr] %s", err)
// 	}
// 	cert, non1, orderloc1, err1 := orresp.finalize(strings.NewReader(string(cout)))
// 	if err1 != nil {
// 		t.Errorf("[finalize] %s", err1)
// 	}
// 	log.Println("non1", non1, "orderloc1", orderloc1, "err1", err1)
// 	err = saveAsFile("cert.pem", cert)
// 	if err != nil {
// 		t.Fatalf("[save cert] %s", err)
// 	}
// 	log.Println("Certificate saved as cert.pem")
// 	err = SavePrivateKeyAsPEM("cert.key", kyss.private)
// 	if err != nil {
// 		t.Fatalf("[save key] %s", err)
// 	}
// 	log.Println("Private key saved as cert.key")

// }
