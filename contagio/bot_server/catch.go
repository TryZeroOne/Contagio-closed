package bot_server

import (
	"contagio/contagio/config/logging"
)

func catch() {
	if err := recover(); err != nil {
		logging.PrintError("Fatal error: " + err.(string))
	}
}
