package main

import (
	"encoding/base64"
	"tvpn/ovpn"
	"fmt"
	"math/big"
)

func main() {
	bobParams := make([]ovpn.Params,4)
	aliceParams := make([]ovpn.Params,4)
	aliceSecret := make([]*big.Int,4)
	bobSecret := make([]*big.Int,4)
	for i := 0; i < 4; i++ {
		bobParams[i] = ovpn.GenParams()
		aliceParams[i] = ovpn.GenParams()
	}
	for i:= 0; i < 4; i++ {
		fmt.Println(base64.StdEncoding.EncodeToString([]byte(bobParams[i].X.Bytes())))
		bobSecret[i] = ovpn.GenMutSecret(bobParams[i],aliceParams[i])
		aliceSecret[i] = ovpn.GenMutSecret(aliceParams[i],bobParams[i])
	}
	bobOut := ovpn.EncodeOpenVPNKey(bobSecret...)
	aliceOut := ovpn.EncodeOpenVPNKey(aliceSecret...)
	fmt.Printf("%s",string(bobOut))
	fmt.Printf("%s",string(aliceOut))
}
