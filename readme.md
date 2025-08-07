# le-cert-core

[![Go](https://github.com/mugund10/le-cert-core/actions/workflows/go.yml/badge.svg)](https://github.com/mugund10/le-cert-core/actions/workflows/go.yml)

`le-cert-core` is a simple ACME (Let's Encrypt) client library written in pure Go.

This project began as a personal journey to deeply understand the ACME protocol, as defined in [RFC 8555](https://datatracker.ietf.org/doc/html/rfc8555). I had previously developed an [ACME client](https://github.com/mugund10/LetsEncryptAcmeClient) using the official [golang.org/x/crypto/acme](https://pkg.go.dev/golang.org/x/crypto/acme) package, but wanted to explore the protocol at a lower level. To achieve that, I built this implementation completely from scratch — including custom JWS (JSON Web Signature) handling — without relying on any third-party libraries or existing ACME implementations.

You can find a sample client in the example/ folder — the main.go file demonstrates how to use this library to register an account, complete a challenge, and request a certificate.

## License

This project is open source and available under the [MIT License](./LICENSE).
