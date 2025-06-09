package main

import (
	"net"
	"os"
)

type NetworkInfo struct {
	Hostname   string   `json:"hostname"`
	IPs        []string `json:"ips"`
	MACs       []string `json:"macs"`
	Interfaces []string `json:"interfaces"`
}

func (a *App) GetNetworkInfo() NetworkInfo {
	var ips, macs, ifaces []string
	hostname, _ := os.Hostname()
	netIfs, _ := net.Interfaces()
	for _, iface := range netIfs {
		ifaces = append(ifaces, iface.Name)
		macs = append(macs, iface.HardwareAddr.String())
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil || ip.To16() != nil {
				ips = append(ips, ip.String())
			}
		}
	}
	return NetworkInfo{
		Hostname:   hostname,
		IPs:        ips,
		MACs:       macs,
		Interfaces: ifaces,
	}
}
