package contract

import (
	"net"
)

type NetAddrStatus string

const Unknown NetAddrStatus = "unknown"

const Whitelisted NetAddrStatus = "whitelisted"

const Blacklisted NetAddrStatus = "blacklisted"

type NetAddrStorage interface {
	BlackListAdd(netAddr string) error
	BlackListRemove(netAddr string) error
	WhiteListAdd(netAddr string) error
	WhiteListRemove(netAddr string) error
	GetStatus(ip net.IP) (NetAddrStatus, error)
	ClearLists() error
}
