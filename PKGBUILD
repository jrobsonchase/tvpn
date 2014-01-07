# $Id: PKGBUILD 199999 2013-11-21 01:21:02Z allan $
# Maintainer: Allan McRae <allan@archlinux.org>
# Contributor: Andreas Radke <andyrtr@archlinux.org>

pkgname=tvpn-git
pkgver=0.r43.g93995e6
pkgrel=1
pkgdesc="Traverse VPN: Peer-to-Peer VPN solution"
arch=('i686' 'x86_64')
url="http://github.com/Pursuit92/tvpn"
license=('GPL2')
depends=('openvpn')
makedepends=('go' 'git')
source=('git://github.com/Pursuit92/tvpn')

pkgver() {
	cd tvpn

	if GITTAG="$(git describe --abbrev=0 --tags 2>/dev/null)"; then
		echo "$(sed -e "s/^${pkgname%%-git}//" -e 's/^[-_/a-zA-Z]\+//' -e 's/[-_+]/./g' <<< ${GITTAG}).r$(git rev-list --count ${GITTAG}..).g$(git log -1 --format="%h")"
	else
		echo "0.r$(git rev-list --count master).g$(git log -1 --format="%h")"
	fi 
}

build() {
	cd ${srcdir}

	mkdir -p Go/src/github.com/Pursuit92/

	local GOPATH=${srcdir}/Go

	ln -s -f ${srcdir}/tvpn ${srcdir}/Go/src/github.com/Pursuit92/

	go get github.com/Pursuit92/tvpn/tvpn
}

package() {
	mkdir -p ${pkgdir}/usr/share/tvpn ${pkgdir}/usr/bin ${pkgdir}/usr/lib/systemd/system

	cp ${srcdir}/Go/bin/tvpn ${pkgdir}/usr/bin/
	cp ${srcdir}/Go/src/github.com/Pursuit92/tvpn/tvpn.config ${pkgdir}/usr/share/tvpn/
	cp ${srcdir}/Go/src/github.com/Pursuit92/tvpn/tvpn.service ${pkgdir}/usr/lib/systemd/system/
}

sha256sums=('SKIP')
