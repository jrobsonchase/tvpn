package ovpn

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
)

const header string = "-----BEGIN OpenVPN Static key V1-----"
const footer string = "-----END OpenVPN Static key V1-----"

func EncodeOpenVPNKey(secrets [][]byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", header)
	for _, v := range secrets {
		i := big.NewInt(0)
		i.SetBytes(v)
		fmt.Fprintf(&buf, "%s\n", hex.EncodeToString(i.Bytes()[:64]))
	}
	fmt.Fprintf(&buf, "%s\n", footer)
	return buf.Bytes()
}
