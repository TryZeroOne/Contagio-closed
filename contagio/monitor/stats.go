package monitor

import (
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func getCpuUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logging.PrintWarning("Can't get cpu usage: " + err.Error())
		return 0
	}
	if int(percent[0]) == 0 {
		percent[0] = 1
	}

	CpuUsage = int(percent[0])
	return percent[0]
}

func getMemoryUsage() int {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	MemoryUsage = int(math.Ceil(memory.UsedPercent))
	return int(math.Ceil(memory.UsedPercent))
}

func isStarted(c *config.Config) (bool, string) {
	ping, err := net.Dial("tcp", c.CncServer)
	if err != nil {
		return false, config.ParseColors("{red}Offline{reset}")
	}
	ping.Close()

	return true, config.ParseColors("{green}Online{reset}")

}

func parseCpu() string {
	var cpu string

	if CpuUsage > 80 {
		cpu = config.ParseColors("{red}" + strconv.Itoa(CpuUsage) + "%{reset}")
		return cpu
	}

	if CpuUsage >= 60 {
		cpu = config.ParseColors("\x1b[38;5;208m" + strconv.Itoa(CpuUsage) + "%{reset}")
		return cpu
	}

	if CpuUsage <= 59 {
		cpu = config.ParseColors("{green}" + strconv.Itoa(CpuUsage) + "%{reset}")
		return cpu

	}

	return cpu

}

func parseMemory() string {
	var mem string

	if MemoryUsage > 80 {
		mem = config.ParseColors("{red}" + strconv.Itoa(MemoryUsage) + "%{reset}")
		return mem
	}

	if MemoryUsage >= 60 {
		mem = config.ParseColors("\x1b[38;5;208m" + strconv.Itoa(MemoryUsage) + "%{reset}")
		return mem
	}

	if MemoryUsage <= 59 {
		mem = config.ParseColors("{green}" + strconv.Itoa(MemoryUsage) + "%{reset}")
		return mem

	}

	return mem

}

func parseTraffic(tr float64) string {
	var result string

	if tr < 40 {
		result = config.ParseColors(fmt.Sprintf("{green}%.2f{reset} KB/s", tr))
		return result
	}

	if tr > 40 {
		result = config.ParseColors(fmt.Sprintf("\x1b[38;5;208m%.2f{reset} KB/s", tr))
		return result
	}

	if tr > 80 {
		result = config.ParseColors(fmt.Sprintf("\x1b[38;5;208m%.2f{reset} KB/s", tr))
		return result
	}

	return result

}
