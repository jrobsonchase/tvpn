package dh

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type Params struct {
	Priv []byte
	X    *big.Int
	Y    *big.Int
}

var curve elliptic.Curve

func init() {
	curve = elliptic.P521()
}

func GenParams() Params {
	priv, x, y, _ := elliptic.GenerateKey(curve, rand.Reader)
	return Params{priv, x, y}
}

func GenMutSecret(local, remote Params) *big.Int {
	secret, _ := curve.ScalarMult(remote.X, remote.Y, local.Priv)
	return secret
}
