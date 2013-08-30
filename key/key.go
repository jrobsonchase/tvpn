package key

import (
	"math/big"
	"encoding/hex"
	"bytes"
	"fmt"
)

const header string = "-----BEGIN OpenVPN Static key V1-----"
const footer string = "-----END OpenVPN Static key V1-----"


func EncodeOpenVPNKey(secret ...*big.Int) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf,"%s\n",header)
	for _,v := range secret {
		fmt.Fprintf(&buf,"%s\n",hex.EncodeToString(v.Bytes()))
	}
	fmt.Fprintf(&buf,"%s\n",footer)
	return buf.Bytes()
}

