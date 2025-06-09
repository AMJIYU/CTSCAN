package main

import (
	"net"
	"os"

	"fmt"
	"strconv"

	gopsnet "github.com/shirou/gopsutil/v4/net"
)

type NetworkInfo struct {
	Hostname   string   `json:"hostname"`
	IPs        []string `json:"ips"`
	MACs       []string `json:"macs"`
	Interfaces []string `json:"interfaces"`
}

type NetworkConn struct {
	Proto      string `json:"proto"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
	Status     string `json:"status"`
	Pid        int32  `json:"pid"`
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

func (a *App) GetNetworkConnections() []NetworkConn {
	conns, err := gopsnet.Connections("all")
	if err != nil {
		return nil
	}
	var result []NetworkConn
	for _, c := range conns {
		proto := fmt.Sprintf("%d", c.Type)
		local := c.Laddr.IP + ":" + strconv.Itoa(int(c.Laddr.Port))
		remote := c.Raddr.IP + ":" + strconv.Itoa(int(c.Raddr.Port))
		result = append(result, NetworkConn{
			Proto:      proto,
			LocalAddr:  local,
			RemoteAddr: remote,
			Status:     c.Status,
			Pid:        c.Pid,
		})
	}
	return result
}
