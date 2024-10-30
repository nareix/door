package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func runListen() error {
	// 这里的 192.168.11.192 是你的 IP 地址，需要替换成你自己的
	l, err := net.Listen("tcp", "192.168.11.192:18022")
	if err != nil {
		return err
	}

	l2addr, _ := net.ResolveUDPAddr("udp", "192.168.255.255:6672")
	l2, err := net.ListenUDP("udp", l2addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			b := make([]byte, 1024)
			n, addr, err := l2.ReadFrom(b)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if n == 21 &&
				b[0] == 0x0 && b[1] == 0x11 && b[2] == 0x01 && b[3] == 0x34 && b[4] == 0x01 {
				// 这里的 11 01 34 01 是门牌号，需要替换成你自己的，以及下面的数据里出现的 11 01 34 01 都要替换
				addr1 := *addr.(*net.UDPAddr)
				addr1.Port = 6672
				log.Println("<-", addr, n)
				l2.WriteTo([]byte{
					0x01, 0x11, 0x01, 0x34, 0x01, 0xc0, 0xa8, 0x0b,
					0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x40, 0x00, 0x00,
				}, &addr1)
			}
		}
	}()

	for {
		c, err := l.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		go func(c0 net.Conn) {
			log.Println("<-", c0.RemoteAddr())

			remoteIP := c0.RemoteAddr().(*net.TCPAddr).IP.String()
			req0 := make([]byte, 36)
			io.ReadFull(c0, req0)
			res0 := []byte{
				0x07, 0xb8, 0x18, 0x00, 0x00, 0x00, 0x72, 0x65,
				0x71, 0x3d, 0x37, 0x30, 0x35, 0x26, 0x71, 0x75,
				0x65, 0x72, 0x79, 0x2a, 0x00, 0x00, 0x01, 0x00,
				0x00, 0x00, 0x03, 0x00, 0x00, 0x01,
			}
			c0.Write(res0)
			c0.Close()

			send := func(b []byte) {
				log.Println(">", remoteIP+":18022")
				c, err := net.Dial("tcp", remoteIP+":18022")
				if err != nil {
					return
				}
				defer c.Close()
				c.Write(b)
				res := make([]byte, 30)
				io.ReadFull(c, res)
			}

			time.Sleep(time.Second * 3)
			send([]byte{
				0x07, 0xb8, 0x1e, 0x00, 0x00, 0x00, 0x72, 0x65,
				0x71, 0x3d, 0x37, 0x31, 0x30, 0x26, 0x71, 0x75,
				0x65, 0x72, 0x79, 0x3d, 0x11, 0x01, 0x00, 0x02,
				0x11, 0x01, 0x34, 0x01, 0x01, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x00,
			})

			time.Sleep(time.Second * 3)
			send([]byte{
				0x07, 0xb8, 0x18, 0x00, 0x00, 0x00, 0x72, 0x65,
				0x71, 0x3d, 0x35, 0x31, 0x38, 0x26, 0x71, 0x75,
				0x65, 0x72, 0x79, 0x3d, 0x22, 0x11, 0x01, 0x00,
				0x02, 0x11, 0x01, 0x34, 0x01, 0x7d,
			})

			time.Sleep(time.Second * 3)
			send([]byte{
				0x07, 0xb8, 0x1e, 0x00, 0x00, 0x00, 0x72, 0x65,
				0x71, 0x3d, 0x37, 0x30, 0x38, 0x26, 0x71, 0x75,
				0x65, 0x72, 0x79, 0x3d, 0x11, 0x01, 0x00, 0x02,
				0x11, 0x01, 0x34, 0x01, 0x01, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x00,
			})
		}(c)
	}
}

type IpEntry struct {
	IfName string `json:"ifname"`
}

const usbIntfName = "enxe04e7a9631bb"

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

func intfLoop() {
	for {
		if usbIntfExists() {
			exec.Command("bash", "-c", "./setusb.sh")
			log.Println("setusb")
		}
		time.Sleep(time.Second * 3)
	}
}

func run() error {
	checkusb := flag.Bool("checkusb", false, "check usb")
	flag.Parse()
	if *checkusb {
		intfLoop()
	}
	return runListen()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
