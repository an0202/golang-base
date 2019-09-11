package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"
	"time"
)

// PortScanner Object
type PortScanner struct {
	ipAddr         string
	startPort      int32
	endPort        int32
	availablePorts []int32
}

// Check Configuration
func (ps *PortScanner) Check(host string, startPort, endPort int32) {
	fmt.Println("[INFO] Configure Checking  ...")
	match, _ := regexp.MatchString("((?:(?:25[0-5]|2[0-4]\\d|(?:1\\d{2}|[1-9]?\\d))\\.){3}(?:25[0-5]|2[0-4]\\d|(?:1\\d{2}|[1-9]?\\d)))", host)
	if match == false {
		fmt.Printf("[ERROR] %s Is Not A Ip Address\n", host)
		os.Exit(2)
	}
	if startPort > endPort {
		fmt.Println("[ERROR] endPort Must Greater Than Or Equal To startPort")
		os.Exit(2)
	}
	fmt.Println("[INFO] Configuration Check Successfully")
}

// Scan Func
func (ps *PortScanner) Scan(host string, startPort, endPort int32) {
	wg := sync.WaitGroup{}
	fmt.Printf("[INFO] Scanning Host %s Port %d To %d ...\n", host, startPort, endPort)
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int32) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
			if err != nil {
				// fmt.Println("[ERROR]", err)
				return
			}
			defer conn.Close()
			// append availableport
			lock := new(sync.Mutex)
			lock.Lock()
			ps.availablePorts = append(ps.availablePorts, port)
			lock.Unlock()
		}(port)
	}
	wg.Wait()
	fmt.Printf("[INFO] Scanning Is Over And %d Available Ports Detected\n", len(ps.availablePorts))
	fmt.Printf("[INFO] Available Ports Is : %d\n", ps.availablePorts)
}

func main() {
	fmt.Println("[INFO] Task Start")
	ps := new(PortScanner)
	ps.ipAddr = "127.0.0.1"
	ps.startPort = 1000
	ps.endPort = 2000
	ps.Check(ps.ipAddr, ps.startPort, ps.endPort)
	ps.Scan(ps.ipAddr, ps.startPort, ps.endPort)
	fmt.Println("[INFO] Task Done")
}

// Todo list:
// 排序实现，端口对应实现，命令行参数添加，时间统计，日志
