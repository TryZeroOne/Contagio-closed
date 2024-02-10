package monitor

import (
	"contagio/contagio/config"
	"sync"
	"time"

	keyboard "github.com/eiannone/keyboard"
)

var (
	CpuUsage    int
	MemoryUsage int
)

// actions
var (
	ACTION_OPEN_CMDLINE = make(chan bool)
)

var (
	lock bool
	wg   sync.WaitGroup
)

func StartMonitor(c *config.Config) {

	go func() {
		for {
			go getMemoryUsage()
			getCpuUsage()
		}

	}()
	go OpenKeyboard(c)
	Clear()

	go func() {
		for {
			<-ACTION_OPEN_CMDLINE
			lock = true
			keyboard.Close()
			Clear()

			PrintTable(c)
			for {
				if CmdLine(c) {
					Clear()
					go OpenKeyboard(c)
					lock = false
					break
				}
			}
		}
	}()

	for {
		table(c)
	}

}

func table(c *config.Config) {

	for {
		if lock {
			break
		}
		wg.Add(1)
		PrintTable(c)
		time.Sleep(time.Duration(c.Monitor.UpdateTime) * time.Millisecond)
		wg.Done()
		wg.Wait()
		if lock {
			break
		}
		Clear()
	}

}
