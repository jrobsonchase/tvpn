package key

import (
	"math/big"
	"encoding/hex"
	"bytes"
	"fmt"
)

const header string = "-----BEGIN OpenVPN Static key V1-----"
const footer string = "-----END OpenVPN Static key V1-----"


func EncodeOpenVPNKey(secret *big.Int) []byte {
	hexStr := hex.EncodeToString(secret.Bytes())
	var buf bytes.Buffer
	fmt.Fprintf(&buf,"%s\n%s\n%s\n",header,hexStr,footer)
	return buf.Bytes()
}

