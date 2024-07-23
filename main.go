package main

import (
	"arp/proto/arp"
	"log"
	"time"
)

var (
// bar = progressbar.Default(100)
)

func main() {
	revSender, err := arp.NewRoute()
	if err != nil {
		log.Fatal(err)
	}
	revSender.Receive()
	revSender.Scan()
	time.Sleep(10 * time.Second)
	revSender.Snooping()
	select {}
}
