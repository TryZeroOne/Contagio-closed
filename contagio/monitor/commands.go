package monitor

import (
	"bufio"
	"contagio/contagio/config"
	"contagio/contagio/database/sqlite"
	"fmt"
	"os"
	"strings"
)

type CommandStruct struct {
	Function    func(string, []string)
	Description string
	Name        string
}

var Commands = map[int]CommandStruct{
	0: {
		Name:        "removeuser",
		Description: "Remove user",
		Function:    RemoveUser,
	},
	1: {
		Name:        "adduser",
		Description: "Add new user",
		Function:    AddUser,
	},
}

var success_message = config.ParseColors("{green}Command executed successfully{reset}")

func Help() {
	for k, v := range Commands {
		fmt.Printf("%d: %s -\t%s\n", k+1, v.Name, v.Description)
	}
}

func Handle(command string) {
	for _, i := range Commands {
		if strings.HasPrefix(command, i.Name) {
			i.Function(command, strings.Split(command, " "))
			return
		}
	}
}

func RemoveUser(command string, command_array []string) {
	login := StringPrompt("Login: ")
	sqlite.RemoveUser(login)
	fmt.Println(success_message)
}

func AddUser(command string, command_array []string) {
	login := StringPrompt("Login: ")
	password := StringPrompt("Password: ")
	sqlite.AddUser(login, password)
	fmt.Println(success_message)
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func CmdLine(c *config.Config) bool {

	str := StringPrompt("Enter command: ")

	if str == ":q" || str == "exit" || str == "quit" {
		fmt.Println("Exit")
		return true
	}

	if str == "help" || str == "?" || str == "Help" || str == "commands" {
		Help()
		return false
	}

	Handle(str)

	return false

}
