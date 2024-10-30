//go:build exclude

package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

const usbIntfName = "enxe04e7a9631bb"

type IpEntry struct {
	IfName string `json:"ifname"`
}

func usbIntfExists() bool {
	b, _ := exec.Command("ip", "-j", "a").Output()
	ips := []IpEntry{}
	json.Unmarshal(b, &ips)
	for _, ip := range ips {
		if ip.IfName == "enxe04e7a9631bb" {
			return true
		}
	}
	return false
}

func main() {
	for {
		if usbIntfExists() {
			exec.Command("bash", "-c", "./setusb.sh")
			log.Println("setusb")
		}
		time.Sleep(time.Second * 3)
	}
}
