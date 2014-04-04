package main

func wrapIpv4Dns(data []uint32) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, convertIpv4AddressToString(a))
	}
	return
}
func unwrapIpv4Dns(wrapData []string) (data []uint32) {
	for _, a := range wrapData {
		data = append(data, convertIpv4AddressToUint32(a))
	}
	return
}

func wrapIpv4Addresses(data [][]uint32) (wrapData Ipv4AddressesWrapper) {
	for _, d := range data {
		if len(d) != 3 {
			LOGGER.Error("ipv4 address invalid", d)
			continue
		}
		ipv4Addr := Ipv4AddressWrapper{}
		ipv4Addr.Address = convertIpv4AddressToString(d[0])
		ipv4Addr.Mask = convertIpv4PrefixToNetMask(d[1])
		ipv4Addr.Gateway = convertIpv4AddressToString(d[2])
		wrapData = append(wrapData, ipv4Addr)
	}
	return
}
func unwrapIpv4Addresses(wrapData Ipv4AddressesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Addr := make([]uint32, 3)
		ipv4Addr[0] = convertIpv4AddressToUint32(d.Address)
		ipv4Addr[1] = convertIpv4NetMaskToPrefix(d.Mask)
		ipv4Addr[2] = convertIpv4AddressToUint32(d.Gateway)
		data = append(data, ipv4Addr)
	}
	return
}

func wrapIpv4Routes(data [][]uint32) (wrapData Ipv4RoutesWrapper) {
	for _, d := range data {
		if len(d) != 4 {
			LOGGER.Error("invalid ipv4 route", d)
			continue
		}
		ipv4Route := Ipv4RouteWrapper{}
		ipv4Route.Address = convertIpv4AddressToString(d[0])
		ipv4Route.Mask = convertIpv4PrefixToNetMask(d[1])
		ipv4Route.NextHop = convertIpv4AddressToString(d[2])
		ipv4Route.Metric = d[3]
		wrapData = append(wrapData, ipv4Route)
	}
	return
}
func unwrapIpv4Routes(wrapData Ipv4RoutesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Route := make([]uint32, 4)
		ipv4Route[0] = convertIpv4AddressToUint32(d.Address)
		ipv4Route[1] = convertIpv4NetMaskToPrefix(d.Mask)
		ipv4Route[2] = convertIpv4AddressToUint32(d.NextHop)
		ipv4Route[3] = d.Metric
		data = append(data, ipv4Route)
	}
	return
}

func wrapIpv6Dns(data [][]byte) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, convertIpv6AddressToString(a))
	}
	return
}
func unwrapIpv6Dns(wrapData []string) (data [][]byte) {
	for _, a := range wrapData {
		data = append(data, convertIpv6AddressToArrayByte(a))
	}
	return
}

func wrapIpv6Addresses(data Ipv6Addresses) (wrapData Ipv6AddressesWrapper) {
	for _, d := range data {
		ipv6Addr := Ipv6AddressWrapper{}
		ipv6Addr.Address = convertIpv6AddressToString(d.Address)
		ipv6Addr.Prefix = d.Prefix
		ipv6Addr.Gateway = convertIpv6AddressToString(d.Gateway)
		wrapData = append(wrapData, ipv6Addr)
	}
	return
}
func unwrapIpv6Addresses(wrapData Ipv6AddressesWrapper) (data Ipv6Addresses) {
	for _, d := range wrapData {
		ipv6Addr := Ipv6Address{}
		ipv6Addr.Address = convertIpv6AddressToArrayByte(d.Address)
		ipv6Addr.Prefix = d.Prefix
		ipv6Addr.Gateway = convertIpv6AddressToArrayByte(d.Gateway)
		data = append(data, ipv6Addr)
	}
	return
}

func wrapIpv6Routes(data Ipv6Routes) (wrapData Ipv6RoutesWrapper) {
	for _, d := range data {
		ipv6Route := Ipv6RouteWrapper{}
		ipv6Route.Address = convertIpv6AddressToString(d.Address)
		ipv6Route.Prefix = d.Prefix
		ipv6Route.NextHop = convertIpv6AddressToString(d.NextHop)
		ipv6Route.Metric = d.Metric
		wrapData = append(wrapData, ipv6Route)
	}
	LOGGER.Debug("ipv6 routes", data, wrapData)
	return
}
func unwrapIpv6Routes(wrapData Ipv6RoutesWrapper) (data Ipv6Routes) {
	for _, d := range wrapData {
		ipv6Route := Ipv6Route{}
		ipv6Route.Address = convertIpv6AddressToArrayByte(d.Address)
		ipv6Route.Prefix = d.Prefix
		ipv6Route.NextHop = convertIpv6AddressToArrayByte(d.NextHop)
		ipv6Route.Metric = d.Metric
		data = append(data, ipv6Route)
	}
	LOGGER.Debug("ipv6 routes", data, wrapData)
	return
}
