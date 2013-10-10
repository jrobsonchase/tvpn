package ovpn

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
)

const header string = "-----BEGIN OpenVPN Static key V1-----"
const footer string = "-----END OpenVPN Static key V1-----"

func EncodeOpenVPNKey(secrets []*big.Int) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", header)
	for _, v := range secrets {
		fmt.Fprintf(&buf, "%s\n", hex.EncodeToString(v.Bytes()[:64]))
	}
	fmt.Fprintf(&buf, "%s\n", footer)
	return buf.Bytes()
}
