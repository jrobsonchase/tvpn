package dh

import (
	"math/big"
)

type Params struct {
	P *big.Int
	G *big.Int
}

func GenPub(secret *big.Int, params Params) *big.Int {
	result := big.NewInt(0)
	result.Exp(params.G,secret,params.P)
	return result
}

func GenMutSecret(secret, remote *big.Int, params Params) *big.Int {
	return big.NewInt(0).Exp(remote,secret,params.P)
}
