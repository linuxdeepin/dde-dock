package main

import "fmt"
import "dlib/dbus"
import nm "dbus/org/freedesktop/networkmanager"

type ActiveConnection struct {
	Interface    string
	HWAddress    string
	IPAddress    string
	SubnetMask   string
	RouteAddress string
	Speed        string
}

func parseIP4address(v uint32) string {
	fmt.Println("Parseip:", v)
	return fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func parseDHCP4(path dbus.ObjectPath) (string, string, string) {
	dhcp4, err := nm.NewDHCP4Config(path)
	if err != nil {
		panic(err)
	}
	options := dhcp4.Options.Get()
	route, _ := options["routers"].Value().(string)
	ip, _ := options["ip_address"].Value().(string)
	mask, _ := options["subnet_mask"].Value().(string)
	return ip, mask, route
}
