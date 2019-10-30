package main

import (
	"fmt"
	"golang-base/tools"
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
	tools.InfoLogger.Println("Configure Checking  ...")
	match, _ := regexp.MatchString("((?:(?:25[0-5]|2[0-4]\\d|(?:1\\d{2}|[1-9]?\\d))\\.){3}(?:25[0-5]|2[0-4]\\d|(?:1\\d{2}|[1-9]?\\d)))", host)
	if match == false {
		tools.ErrorLogger.Fatalf("%s Is Not A Ip Address\n", host)
	}
	if startPort > endPort {
		tools.ErrorLogger.Fatalln("endPort Must Greater Than Or Equal To startPort")
	}
	tools.InfoLogger.Println("Configuration Check Successfully")
}

// Scan Func
func (ps *PortScanner) Scan(host string, startPort, endPort int32) {
	wg := sync.WaitGroup{}
	tools.InfoLogger.Printf("Scanning Host %s Port %d To %d ...\n", host, startPort, endPort)
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
	tools.InfoLogger.Printf("%d Available Ports Detected\n", len(ps.availablePorts))
	tools.InfoLogger.Printf("Available Ports Is : %d\n", ps.availablePorts)
}

func (ps *PortScanner) descriebPorts() {
	if len(ps.availablePorts) == 0 {
		os.Exit(1)
	}
	for _, port := range ps.availablePorts {
		result := tools.DescribePort(port)
		if result != "" {
			fmt.Printf("%-5d --> %6s\n", port, result)
		} else {
			fmt.Printf("%-5d --> %6s\n", port, "unknown Port")
		}
	}
}

func main() {
	tools.InfoLogger.Println("Task Start")
	ps := new(PortScanner)
	ps.ipAddr = "180.76.166.25"
	ps.startPort = 1
	ps.endPort = 65535
	ps.Check(ps.ipAddr, ps.startPort, ps.endPort)
	ps.Scan(ps.ipAddr, ps.startPort, ps.endPort)
	ps.descriebPorts()
	tools.InfoLogger.Println("Task Done")
}

// Todo list:
// 排序实现，命令行参数添加
