package main

import (
	"crypto/elliptic"
	"tvpn/key"
	"fmt"
	"crypto/rand"
)

func main() {
	curve := elliptic.P521()
	bobPriv1,_,_,_ := elliptic.GenerateKey(curve,rand.Reader)
	_,aliceX1,aliceY1,_ := elliptic.GenerateKey(curve,rand.Reader)
	bobMut1,_ := curve.ScalarMult(aliceX1,aliceY1,bobPriv1)
	bobPriv2,_,_,_ := elliptic.GenerateKey(curve,rand.Reader)
	_,aliceX2,aliceY2,_ := elliptic.GenerateKey(curve,rand.Reader)
	bobMut2,_ := curve.ScalarMult(aliceX2,aliceY2,bobPriv2)
	bobPriv3,_,_,_ := elliptic.GenerateKey(curve,rand.Reader)
	_,aliceX3,aliceY3,_ := elliptic.GenerateKey(curve,rand.Reader)
	bobMut3,_ := curve.ScalarMult(aliceX3,aliceY3,bobPriv3)
	bobPriv4,_,_,_ := elliptic.GenerateKey(curve,rand.Reader)
	_,aliceX4,aliceY4,_ := elliptic.GenerateKey(curve,rand.Reader)
	bobMut4,_ := curve.ScalarMult(aliceX4,aliceY4,bobPriv4)
	out := key.EncodeOpenVPNKey(bobMut1,bobMut2,bobMut3,bobMut4)
	fmt.Printf("%s",string(out))
}
