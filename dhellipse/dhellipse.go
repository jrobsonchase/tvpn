package main

import (
	"tvpn/key"
	"fmt"
	"tvpn/dh"
	"math/big"
)

func main() {
	bobParams := make([]dh.Params,4)
	aliceParams := make([]dh.Params,4)
	aliceSecret := make([]*big.Int,4)
	bobSecret := make([]*big.Int,4)
	for i := 0; i < 4; i++ {
		bobParams[i] = dh.GenParams()
		aliceParams[i] = dh.GenParams()
	}
	for i:= 0; i < 4; i++ {
		bobSecret[i] = dh.GenMutSecret(bobParams[i],aliceParams[i])
		aliceSecret[i] = dh.GenMutSecret(aliceParams[i],bobParams[i])
	}
	bobOut := key.EncodeOpenVPNKey(bobSecret...)
	aliceOut := key.EncodeOpenVPNKey(aliceSecret...)
	fmt.Printf("%s",string(bobOut))
	fmt.Printf("%s",string(aliceOut))
}
