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

package tvpn

import (
	"math"
	"net"
	"strconv"
	"github.com/Pursuit92/LeveledLogger/log"
)

type IPConfig map[string]string

type IPReq struct {
	Req  bool
	IP   net.IP
	Resp chan net.IP
}

type IPManager struct {
	reqs  chan IPReq
	reinit bool
	Start net.IP
	Tuns  int
}

func (ipman *IPManager) Init() {
	if ipman.reinit {
		ipman.stopCurrent()
	}
	ipman.reqs = make(chan IPReq)
	ipman.reinit = true
	go ipAllocator(ipman.reqs,ipman.Start,ipman.Tuns)
}

func (ipman *IPManager) Configure(conf IPConfig) {
	ipman.Start = net.ParseIP(conf["Start"])
	num,err := strconv.Atoi(conf["Num"])
	if err != nil {
		panic(err)
	}
	ipman.Tuns = num
	ipman.Init()
}

func (ipman IPManager) RequestAny() net.IP {
	resp := make(chan net.IP)
	req := IPReq{Req: true, Resp: resp}
	ipman.reqs <- req
	return <-resp
}

func (ipman IPManager) Request(ip net.IP) net.IP {
	log.Out.Lprintln(3,"Attempting to allocate ",ip)
	resp := make(chan net.IP)
	req := IPReq{Req: true, IP: ip, Resp: resp}
	ipman.reqs <- req
	ret := <-resp
	log.Out.Lprintln(3,"Returning ",ret)
	return ret
}

func (ipman IPManager) Release(ip net.IP) net.IP {
	resp := make(chan net.IP)
	req := IPReq{Req: false, IP: ip, Resp: resp}
	ipman.reqs <- req
	return <-resp
}

func (ipman IPManager) stopCurrent() {
	close(ipman.reqs)
}

func ipAllocator(ipReqs chan IPReq, min net.IP, n int) {
	allocList := make([]bool, n)
	for req := range ipReqs {
		// is it a request for an IP or relinquishing one?
		if req.Req {
			// is it a request for any ip or a specific one?
			if req.IP == nil {
				log.Out.Lprintln(4,"Got request for first available!")
				// any case, pick the first unallocated
				for i, v := range allocList {
					if !v {
						allocList[i] = true
						req.Resp <- indexToIP(min, i)
						break
					}
				}
			} else {
				log.Out.Lprintln(4,"Attempting to allocate ",req.IP)
				// specific: if the requested isn't available, pick the next
				i := ipToIndex(min, req.IP)
				if !allocList[i] {
					log.Out.Lprintln(4,"Allocated ",indexToIP(min,i))
					req.Resp <- indexToIP(min, i)
					allocList[i] = true
				} else {
					for j := i; j < len(allocList); j++ {
						if !allocList[j] {
							log.Out.Lprintln(4,"Allocated ",indexToIP(min,j))
							req.Resp <- indexToIP(min, j)
							allocList[j] = true
							break
						}
					}
				}
			}

		} else {
			//relinquish case
			i := ipToIndex(min, req.IP)
			allocList[i] = false
			req.Resp <- indexToIP(min, i)
		}
	}

}

func ipToIndex(start, ip net.IP) int {
	start4 := start.To4()
	ip4 := ip.To4()
	dif := net.IPv4(ip4[0]-start4[0],
		ip4[1]-start4[1],
		ip4[2]-start4[2],
		ip4[3]-start4[3])
	dif4 := dif.To4()
	var sum int
	for i, v := range dif4 {
		sum += int(float64(v) * math.Pow(256, float64(3-i)))
	}
	return sum / 4
}

func indexToIP(start net.IP, index int) net.IP {
	index *= 4
	bs := make([]byte, 4)
	for i, v := range start.To4() {
		bs[i] = v + byte(float64(index)/math.Pow(256, float64(3-i)))
	}
	return net.IPv4(bs[0], bs[1], bs[2], bs[3])
}

func isGreater(lhs,rhs net.IP) bool {
	lhs4 := lhs.To4()
	rhs4 := rhs.To4()
	for i := range lhs4 {
		if lhs4[i] > rhs4[i] {
			return true
		}
	}
	return false
}
