package arp

import (
	"arp/model"
	"arp/utils"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/google/gopacket/pcap"
	"net"
	"time"
)

// 获取IP地址范围
func getIPRange(local pcap.InterfaceAddress) []net.IP {
	var ips []net.IP
	start := utils.IPToUInt32(local.IP.Mask(local.Netmask))
	ones, bits := local.Netmask.Size()
	times := utils.Pow(2, uint32(bits-ones))
	var i uint32
	for i = 0; i < times; i++ {
		tmp := utils.UInt32ToIP(start + i)
		if tmp != nil {
			ips = append(ips, tmp)
		}
	}
	return ips
}

func (r *Route) Scan() {
	IPs := getIPRange(r.srcIPAddr)
	for _, ip := range IPs {
		ip := ip
		r.wg.Add(1)
		go func() {
			pinger, err := ping.NewPinger(ip.String())
			pinger.SetPrivileged(true)
			pinger.Timeout = 5 * time.Second
			if err != nil {
				fmt.Println("err")
				r.wg.Done()
				return
			}
			pinger.Count = 3
			err = pinger.Run() // Blocks until finished.
			if err == nil && !r.srcIPAddr.IP.Equal(ip) {
				stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
				err := r.SendARP(r.gateWayIP, ip, model.Request, r.srcMAC, BroadcastMac)
				if err != nil {
					fmt.Errorf("find IP:%s,but occurred error:%e", stats.Addr, err)
				}
			}
			r.wg.Done()
		}()
	}
	return
}
