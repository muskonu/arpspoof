package utils

import (
	"github.com/jackpal/gateway"
	"log"
	"net"
)

func GetGateway() net.IP {
	ip, err := gateway.DiscoverGateway()
	if err != nil {
		log.Fatal(err)
	}
	return ip
}
