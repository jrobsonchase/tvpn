package main

import (
	"github.com/Pursuit92/tvpn"
	"github.com/codegangsta/martini"
)

func createAPI(tinst *tvpn.TVPN) (m *martini.ClassicMartini) {
	m = martini.Classic()
	m.Get(    "/api/v1.0/", func() {})
	m.Post(   "/api/v1.0/", func() {
		//tinst.Stop()
		tinst.Start()
	})
	m.Delete( "/api/v1.0/", func() {
		tinst.Stop()
		tinst.Cleanup()
	})

	m.Get(    "/api/v1.0/signaling/", func() {})
	m.Put(    "/api/v1.0/signaling/", func() {})
	m.Post(   "/api/v1.0/signaling/", func() {})
	m.Delete( "/api/v1.0/signaling/", func() {})

	m.Get(    "/api/v1.0/signaling/name/", func() {})
	m.Put(    "/api/v1.0/signaling/name/", func() {})

	m.Get(    "/api/v1.0/signaling/group", func() {})
	m.Put(    "/api/v1.0/signaling/group", func() {})

	m.Get(    "/api/v1.0/signaling/server/", func() {})
	m.Put(    "/api/v1.0/signaling/server/", func() {})

	m.Get(    "/api/v1.0/stun/", func() {})
	m.Put(    "/api/v1.0/stun/", func() {})

	m.Get(    "/api/v1.0/ipalloc/", func() {})
	m.Put(    "/api/v1.0/ipalloc/", func() {})

	m.Get(    "/api/v1.0/ipalloc/tunnels/", func() {})

	m.Get(    "/api/v1.0/ipalloc/max/", func() {})
	m.Put(    "/api/v1.0/ipalloc/max/", func() {})

	m.Get(    "/api/v1.0/ipalloc/base/", func() {})
	m.Put(    "/api/v1.0/ipalloc/base/", func() {})

	m.Get(    "/api/v1.0/vpn/", func() {})
	m.Put(    "/api/v1.0/vpn/", func() {})

	m.Get(    "/api/v1.0/vpn/tmp/", func() {})
	m.Put(    "/api/v1.0/vpn/tmp/", func() {})

	m.Get(    "/api/v1.0/vpn/path/", func() {})
	m.Put(    "/api/v1.0/vpn/path/", func() {})

	m.Get(    "/api/v1.0/friends/", func() {})
	m.Post(   "/api/v1.0/friends/", func() {})
	m.Delete( "/api/v1.0/friends/", func() {})

	m.Get(    "/api/v1.0/friends/:name", func() {})
	m.Put(    "/api/v1.0/friends/:name", func() {})
	m.Delete( "/api/v1.0/friends/:name", func() {})

	m.Get(    "/api/v1.0/friends/:name/validate/", func() {})
	m.Put(    "/api/v1.0/friends/:name/validate/", func() {})

	m.Get(    "/api/v1.0/friends/:name/routes/", func() {})
	m.Put(    "/api/v1.0/friends/:name/routes/", func() {})

	m.Get(    "/api/v1.0/friends/:name/state/", func() {})
	m.Post(   "/api/v1.0/friends/:name/state/", func() {})
	m.Delete( "/api/v1.0/friends/:name/state/", func() {})

	m.Get(    "/api/v1.0/friends/:name/log/", func() {})
	m.Delete( "/api/v1.0/friends/:name/log/", func() {})

	return m
}
