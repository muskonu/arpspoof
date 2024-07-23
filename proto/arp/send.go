package arp

import (
	"arp/model"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

func (r *Route) SendARP(srcIP, dstIP net.IP, operation model.Operation, srcMAC net.HardwareAddr, dstMAC net.HardwareAddr) error {
	//构造ARP数据包
	eth := layers.Ethernet{
		SrcMAC:       r.srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         uint16(operation), //request or reply
		SourceHwAddress:   srcMAC,
		SourceProtAddress: srcIP,
		DstHwAddress:      dstMAC,
		DstProtAddress:    dstIP.To4(),
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		fmt.Println(arp)
		log.Println("gopacket.SerializeLayers failed,err:", err)
		return err
	}
	//发送ARP数据包
	err := r.handle.WritePacketData(buf.Bytes())
	if err != nil {
		fmt.Println("handle.WritePacketData failed,err:", err)
		return err
	}
	return nil
}
