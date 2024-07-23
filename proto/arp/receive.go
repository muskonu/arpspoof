package arp

import (
	"arp/model"
	"bytes"
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	manuf "github.com/timest/gomanuf"
	"log"
	"net"
)

func (r *Route) Receive() {
	go r.arpReceive(r.ctx)
	return
}

func (r *Route) arpReceive(ctx context.Context) {
	log.Println("开始接收ARP报文")
	defer func() {
		log.Println("停止接收")
	}()
	//1.创建数据源
	src := gopacket.NewPacketSource(r.handle, layers.LayerTypeEthernet)
	//2.从数据源中读数据
	for packet := range src.Packets() {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil { //不是ARP包
			continue
		}
		arp, ok := arpLayer.(*layers.ARP)
		if !ok {
			continue
		}
		//不是ARP响应包
		if arp.Operation != layers.ARPReply {
			continue
		}
		//不是发送给我的包
		if false == bytes.Equal(r.srcMAC, arp.DstHwAddress) {
			continue
		}
		if r.gateWayIP.Equal(arp.SourceProtAddress) {
			r.gateWayMac = arp.SourceHwAddress
			continue
		}
		host := &model.Host{
			IP:      net.IP(arp.SourceProtAddress).String(),
			MAC:     net.HardwareAddr(arp.SourceHwAddress).String(),
			MACInfo: manuf.Search(net.HardwareAddr(arp.SourceHwAddress).String()),
		}
		r.mutex.Lock()
		if _, exist := r.Hosts[host.MAC]; !exist {
			fmt.Println(host)
			r.Hosts[host.MAC] = host
		}
		r.mutex.Unlock()
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
