package pkg

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	gopsnet "github.com/shirou/gopsutil/v4/net"
)

type InterfaceStats struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

type NetworkInfo struct {
	Hostname       string           `json:"hostname"`
	IPs            []string         `json:"ips"`
	MACs           []string         `json:"macs"`
	Interfaces     []string         `json:"interfaces"`
	InterfaceStats []InterfaceStats `json:"interface_stats"`
	Gateway        string           `json:"gateway"`
}

type NetworkConn struct {
	Proto      string `json:"proto"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
	Status     string `json:"status"`
	Pid        int32  `json:"pid"`
}

type RouteInfo struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Interface   string `json:"interface"`
	Flags       string `json:"flags"`
}

func (a *App) GetNetworkInfo() NetworkInfo {
	var ips, macs, ifaces []string
	hostname, _ := os.Hostname()
	netIfs, _ := net.Interfaces()
	gateway := getDefaultGateway()

	// 获取网络流量统计
	ioStats, _ := gopsnet.IOCounters(true)
	interfaceStats := make([]InterfaceStats, 0)

	for _, iface := range netIfs {
		ifaces = append(ifaces, iface.Name)
		mac := iface.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		} else {
			macs = append(macs, "")
		}

		// 查找对应网卡的流量统计
		for _, stat := range ioStats {
			if stat.Name == iface.Name {
				interfaceStats = append(interfaceStats, InterfaceStats{
					Name:        iface.Name,
					BytesSent:   stat.BytesSent,
					BytesRecv:   stat.BytesRecv,
					PacketsSent: stat.PacketsSent,
					PacketsRecv: stat.PacketsRecv,
				})
				break
			}
		}

		addrs, _ := iface.Addrs()
		var interfaceIPs []string
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
				interfaceIPs = append(interfaceIPs, ip.String())
			}
		}
		if len(interfaceIPs) > 0 {
			ips = append(ips, interfaceIPs[0]) // 只取第一个 IP
		} else {
			ips = append(ips, "")
		}
	}

	return NetworkInfo{
		Hostname:       hostname,
		IPs:            ips,
		MACs:           macs,
		Interfaces:     ifaces,
		InterfaceStats: interfaceStats,
		Gateway:        gateway,
	}
}

func getDefaultGateway() string {
	switch runtime.GOOS {
	case "windows":
		output, err := backgroundCommand("route.exe", "print", "-4", "0.0.0.0").Output()
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(output), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 5 && fields[0] == "0.0.0.0" && fields[1] == "0.0.0.0" && fields[2] != "On-link" {
				return fields[2]
			}
		}
	case "darwin":
		output, err := backgroundCommand("route", "-n", "get", "default").Output()
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "gateway") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					return fields[len(fields)-1]
				}
			}
		}
	case "linux":
		output, err := backgroundCommand("ip", "route", "show", "default").Output()
		if err != nil {
			return ""
		}
		fields := strings.Fields(string(output))
		for i, field := range fields {
			if field == "via" && i+1 < len(fields) {
				return fields[i+1]
			}
		}
	}

	return ""
}

func protoName(t uint32) string {
	switch t {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	default:
		return fmt.Sprintf("%d", t)
	}
}

func (a *App) GetNetworkConnections() []NetworkConn {
	conns, err := gopsnet.Connections("all")
	if err != nil {
		return nil
	}
	var result []NetworkConn
	for _, c := range conns {
		proto := protoName(c.Type)
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
