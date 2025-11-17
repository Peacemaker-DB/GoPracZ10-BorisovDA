package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func save(filename string, b []byte) {
	os.WriteFile(filename, b, 0644)
}

func main() {
	key1, _ := rsa.GenerateKey(rand.Reader, 2048)
	priv1 := x509.MarshalPKCS1PrivateKey(key1)
	pub1 := x509.MarshalPKCS1PublicKey(&key1.PublicKey)

	save("private1.pem", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: priv1}))
	save("public1.pem", pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pub1}))

	key2, _ := rsa.GenerateKey(rand.Reader, 2048)
	priv2 := x509.MarshalPKCS1PrivateKey(key2)
	pub2 := x509.MarshalPKCS1PublicKey(&key2.PublicKey)

	save("private2.pem", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: priv2}))
	save("public2.pem", pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pub2}))

	println("Готово: private1.pem, public1.pem, private2.pem, public2.pem")
}
