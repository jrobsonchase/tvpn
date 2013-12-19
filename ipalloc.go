package tvpn

import (
	"fmt"
	"math"
	"net"
)

type IPReq struct {
	Req  bool
	IP   net.IP
	Resp chan net.IP
}

type IPManager struct {
	reqs  chan IPReq
	Start net.IP
	Tuns  int
}

func NewIPManager(minS string, numTun int) IPManager {
	min := net.ParseIP(minS)
	reqs := make(chan IPReq)
	go ipAllocator(reqs, min, numTun)
	return IPManager{reqs, min, numTun}
}

func (ipman IPManager) RequestAny() net.IP {
	resp := make(chan net.IP)
	req := IPReq{Req: true, Resp: resp}
	ipman.reqs <- req
	return <-resp
}

func (ipman IPManager) Request(ip net.IP) net.IP {
	fmt.Println("Attempting to allocate ",ip)
	resp := make(chan net.IP)
	req := IPReq{Req: true, IP: ip, Resp: resp}
	ipman.reqs <- req
	ret := <-resp
	fmt.Println("Returning ",ret)
	return ret
}

func (ipman IPManager) Release(ip net.IP) net.IP {
	resp := make(chan net.IP)
	req := IPReq{Req: false, IP: ip, Resp: resp}
	ipman.reqs <- req
	return <-resp
}

func ipAllocator(ipReqs chan IPReq, min net.IP, n int) {
	allocList := make([]bool, n)
	for {
		req := <-ipReqs
		// is it a request for an IP or relinquishing one?
		if req.Req {
			// is it a request for any ip or a specific one?
			if req.IP == nil {
				fmt.Println("Got request for first available!")
				// any case, pick the first unallocated
				for i, v := range allocList {
					if !v {
						allocList[i] = true
						req.Resp <- indexToIP(min, i)
						break
					}
				}
			} else {
				fmt.Println("Attempting to allocate ",req.IP)
				// specific: if the requested isn't available, pick the next
				i := ipToIndex(min, req.IP)
				if !allocList[i] {
					fmt.Println("Allocated ",indexToIP(min,i))
					req.Resp <- indexToIP(min, i)
					allocList[i] = true
				} else {
					for j := i; j < len(allocList); j++ {
						if !allocList[j] {
							fmt.Println("Allocated ",indexToIP(min,j))
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
	for i,_ := range lhs4 {
		if lhs4[i] > rhs4[i] {
			return true
		}
	}
	return false
}

func isEqual(lhs,rhs net.IP) bool {
	fmt.Println("Checking IP equality: ",lhs," == ",rhs,"?")
	lhs4 := lhs.To4()
	rhs4 := rhs.To4()
	for i,_ := range lhs4 {
		fmt.Println("Testing ",lhs4[i]," and ",rhs4[i])
		if lhs4[i] != rhs4[i] {
			return false
		}
	}
	return true
}
