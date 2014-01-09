/*
 *  TVPN: A Peer-to-Peer VPN solution for traversing NAT firewalls
 *  Copyright (C) 2013  Joshua Chase <jcjoshuachase@gmail.com>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package ovpn

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const header string = "-----BEGIN OpenVPN Static key V1-----"
const footer string = "-----END OpenVPN Static key V1-----"

func EncodeOpenVPNKey(secrets [][64]byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", header)
	for _, v := range secrets {
		fmt.Fprintf(&buf, "%s\n", hex.EncodeToString(v[:]))
	}
	fmt.Fprintf(&buf, "%s\n", footer)
	return buf.Bytes()
}
