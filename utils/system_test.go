package utils

import (
	"github.com/jackpal/gateway"
	"testing"
)

func TestGetGateway(t *testing.T) {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Gateway:", gw.String())
	device, addr := GetIPAddr(gw)
	t.Log(addr)
	mac, err := GetMac(device.Name)
	if err != nil {
		t.Fatalf("GetInterfaceByName Error:%v", err)
	}
	t.Log(mac)
}
