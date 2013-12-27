Traverse Peer-to-Peer VPN System
================================

Traverse VPN is a peer-to-peer VPN system utilizing public messaging and
STUN servers to allow clients to traverse NAT firewalls and make direct
connections using VPN tunnels

Development is in a constant state of flux, so details will be sparse
until solidified.

The current implementation utilizes IRC networks for connection
negotiation and OpenVPN to create the actual VPN tunnels. The tvpn
client can be run on both Linux and Windows systems that have OpenVPN
installed, though windows users are limited to one connection (unless
they install more TAP OpenVPN drivers)
