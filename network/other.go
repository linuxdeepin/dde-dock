package main

type IPv4Config struct {
	IPAddress        string
	BroadcastAddress string
	SubnetMask       string
	DefaultRoute     string
	PrimaryDNS       string
}
type ActiveConnection struct {
	Interface string
	HWAddress string
	Drive     string
	Speed     string
	Security  string

	IPv4 IPv4Config
	IPv6 IPv4Config
}

