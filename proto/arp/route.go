package arp

import (
	"arp/model"
	"arp/utils"
	"context"
	"fmt"
	"github.com/google/gopacket/pcap"
	"net"
	"sync"
)

var (
	BroadcastMac net.HardwareAddr
)

func init() {
	BroadcastMac, _ = net.ParseMAC("ff:ff:ff:ff:ff:ff")
}

type Route struct {
	handle *pcap.Handle
	wg     sync.WaitGroup
	ctx    context.Context
	mutex  sync.Mutex
	Hosts  map[string]*model.Host
	cancel context.CancelFunc

	srcMAC     net.HardwareAddr
	gateWayIP  net.IP
	gateWayMac net.HardwareAddr
	srcIPAddr  pcap.InterfaceAddress
}

func NewRoute() (*Route, error) {
	gw := utils.GetGateway().To4()
	device, addr := utils.GetIPAddr(gw)
	handle, err := pcap.OpenLive(device.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("OpenLive Error:%v", err)
	}
	mac, err := utils.GetMac(device.Name)
	if err != nil {
		return nil, fmt.Errorf("GetMac Error:%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Route{
		handle: handle,
		ctx:    ctx,
		Hosts:  make(map[string]*model.Host),
		cancel: cancel,

		gateWayIP: gw,
		srcMAC:    mac,
		srcIPAddr: addr,
	}, nil
}
