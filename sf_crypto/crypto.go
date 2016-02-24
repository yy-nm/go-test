// author: mard

// used to crypto
package sf_crypto

import (
	"math/big"
)

const (
	b int64 = 7
)

var mod = big.NewInt(0).SetUint64(0xffffffffffffffc5)
var base = big.NewInt(b)

func DH_exchange(key *big.Int) *big.Int {

	return DH_secret(base, key)
}

func DH_secret(key1, key2 *big.Int) *big.Int {
	result := big.NewInt(0)
	return result.Exp(key1, key2, mod)
}
