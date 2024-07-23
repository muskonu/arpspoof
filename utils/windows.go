//go:build windows

package utils

import (
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
)

func getCommandConn(command string) string {
	//需要执行命令： command
	cmd := exec.Command("cmd", "/C", command)
	// 获取输入
	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("无法获取命令的标准输出管道", err.Error())
		return ""
	}
	// 执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Windows命令执行失败，请检查命令输入是否有误", err.Error())
		return ""
	}
	enc := mahonia.NewDecoder("gbk")
	// 读取输出
	bytes, err := ioutil.ReadAll(output)
	if err != nil {
		fmt.Println("打印异常，请检查")
		return ""
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait", err.Error())
		return ""
	}
	return enc.ConvertString(string(bytes))
}

func GetMac(ifaceName string) (net.HardwareAddr, error) {
	guid := findGuid(ifaceName)
	gm := getCommandConn("getmac")
	linesGateways := strings.Split(gm, "\n")
	for _, line := range linesGateways {
		if strings.Contains(line, guid) {
			files := strings.Fields(line)
			if mac, err := net.ParseMAC(files[0]); err == nil {
				return mac, nil
			}
		}
	}
	return nil, errors.New("can't find mac address")
}

func GetIPAddr(srcIP net.IP) (pcap.Interface, pcap.InterfaceAddress) {
	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// 打印设备信息
	var device pcap.Interface
	var address pcap.InterfaceAddress
	fmt.Println("Devices found:")

	for _, d := range devices {
		for _, addr := range d.Addresses {
			if addr.IP.To4().Mask(addr.Netmask).String() == srcIP.Mask(srcIP.DefaultMask()).String() {
				fmt.Println(srcIP.Mask(srcIP.DefaultMask()).String())
				device = d
				address = addr
			}
		}
	}

	if address.IP == nil {
		log.Fatalf("No IPv4 address")
	}
	fmt.Println("- IP address: ", address.IP)
	fmt.Println("- Subnet mask: ", address.Netmask)
	fmt.Println("- BroadAddr: ", address.Broadaddr)
	return device, address
}

func findGuid(ifaceName string) string {
	i := strings.Index(ifaceName, "{")
	return ifaceName[i : len(ifaceName)-1]
}
