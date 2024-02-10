package cnc

import (
	"contagio/contagio/config/logging"
)

func Catch() {
	if err := recover(); err != nil {
		logging.PrintError("Fatal error: " + err.(string))
		return
	}
}
