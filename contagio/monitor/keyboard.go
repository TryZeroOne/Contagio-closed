package monitor

import (
	"contagio/contagio/config"

	"contagio/contagio/config/logging"
	"os"

	keyboard "github.com/eiannone/keyboard"
)

var Stoplol bool

func OpenKeyboard(c *config.Config) {

	if err := keyboard.Open(); err != nil {
		logging.PrintError("Can't open keyboard: " + err.Error())
		os.Exit(0)
	}

	defer keyboard.Close()

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			logging.PrintWarning("Can't get key: " + err.Error())
			continue
		}

		if key == keyboard.Key(c.Monitor.CloseMonitorKeyUint16) {
			os.Exit(0)
		}
		if key == keyboard.Key(c.Monitor.OpenCmdLineKeyUint16) {
			ACTION_OPEN_CMDLINE <- true
			return
		}

	}

}
