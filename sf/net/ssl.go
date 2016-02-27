package net

type Ssl interface {
	Crypto([]byte) []byte
	Decrypt([]byte) []byte
	//Wait_for_crypto()
}
