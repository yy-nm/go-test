package sf_net

type Ssl interface {
	crypto([]byte) []byte
	decrypt([]byte) []byte
	wait_for_crypto()
}
