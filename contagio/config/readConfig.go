package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
)

var Disable_Debug bool

type Config struct {
	Wg *sync.WaitGroup

	ImportTheme   string
	CncServer     string
	BotServer     string
	TftpServer    string
	FtpServer     string
	LoaderServer  string
	DISABLE_DEBUG bool

	Payload struct {
		FtpPassword string
		FtpLogin    string
	}

	RootLogin string
	Logs      struct {
		TelegramBotToken           string
		TelegramChatId             string
		NewClientConnectedTelegram string
		SaveLogsInFile             bool
		PrintLogsInTerminal        bool
		SendLogsInTelegram         bool

		NewClientConnectedFile     string
		NewClientConnectedFileName string
		NewClientConnectedLog      bool
		NewClientConnectedTerminal string

		NewAttackStartedLog      bool
		NewAttackStartedFileName string
		NewAttackStartedTerminal string
		NewAttackStartedFile     string
		NewAttackStartedTelegram string
	}
	Cnc struct {
		HelpCommand string
		CmdPrompt   string

		MethodsCommand       string
		CustomMethods        []string
		CustomMethodsEnabled bool

		CustomHelp                []string
		CustomHelpEnabled         bool
		Banner                    []string
		InvalidCommandSyntaxError string
		UnknownCommandError       string
		BotCount                  string
		NoBotsConnectedError      string

		NoActiveAttacksError  string
		AttackIdNotFoundError string
		IpNotFoundError       string

		CommandSent string
		Title       string

		CommandExecuted      string
		CommandInvalidSyntax string
		InvalidSubnetError   string
		SubnetCommandSent    string
	}

	Auth struct {
		LoginPrompt         string
		PasswordPrompt      string
		AuthError           string
		CaptchaPrompt       string
		CaptchaError        string
		AllowAllIps         bool
		IpIsNotAllowedError string
	}

	Captcha struct {
		Enabled    bool
		CaptchaLen int
		Letters    string
	}

	Animation struct {
		Enabled bool

		Delay   int
		Letters string
	}

	Logging struct {
		LogLevel       string
		WarningMessage string
		InfoMessage    string
		ErrorMessage   string
	}

	Monitor struct {
		IsCmdLineRunning      chan bool
		OpenCmdLineKeyUint16  uint16
		CloseMonitorKeyUint16 uint16

		///
		UpdateTime      int
		OpenCmdLineKey  string
		CloseMonitorKey string
	}

	Modules map[string]Cmodule
	Ti      ThemeInit
}

type ThemeInit struct {
	ThemeInit struct {
		Author      string /* Optional */
		Version     string
		Description string /* Optional */
		Name        string
	}
}

type Cmodule struct {
	Exec    string
	ExecEnv string
	ExecDir string
}

var Config_ *Config

func ReadConfig(wg *sync.WaitGroup) (*Config, error) {

	var config Config
	var ti ThemeInit

	conf, err := os.ReadFile("./config.toml")
	if err != nil {
		return nil, err
	}
	_, err = toml.Decode(string(conf), &config)

	if err != nil {
		return nil, err
	}

	theme, err := os.ReadFile(config.ImportTheme)
	if err != nil {
		return nil, err
	}
	_, err = toml.Decode(string(theme), &config)
	if err != nil {
		return nil, err
	}

	_, err = toml.Decode(string(theme), &ti) // get Theme init
	if err != nil {
		return nil, err
	}

	if ti.ThemeInit.Version == "" || ti.ThemeInit.Name == "" {
		return nil, fmt.Errorf("theme name and version are required")
	}

	if err := parse(&config); err != nil {
		return nil, err
	}
	Disable_Debug = config.DISABLE_DEBUG
	config.Wg = wg
	config.Ti = ti

	Config_ = &config

	return &config, nil

}

var keylist = map[string]uint16{
	"ctrl+tilde":      0x00,
	"ctrl+2":          0x00,
	"ctrl+space":      0x00,
	"ctrl+a":          0x01,
	"ctrl+b":          0x02,
	"ctrl+c":          0x03,
	"ctrl+d":          0x04,
	"ctrl+e":          0x05,
	"ctrl+f":          0x06,
	"ctrl+g":          0x07,
	"backspace":       0x08,
	"ctrl+h":          0x08,
	"tab":             0x09,
	"ctrl+i":          0x09,
	"ctrl+j":          0x0A,
	"ctrl+k":          0x0B,
	"ctrl+l":          0x0C,
	"enter":           0x0D,
	"ctrl+m":          0x0D,
	"ctrl+n":          0x0E,
	"ctrl+o":          0x0F,
	"ctrl+p":          0x10,
	"ctrl+q":          0x11,
	"ctrl+r":          0x12,
	"ctrl+s":          0x13,
	"ctrl+t":          0x14,
	"ctrl+u":          0x15,
	"ctrl+v":          0x16,
	"ctrl+w":          0x17,
	"ctrl+x":          0x18,
	"ctrl+y":          0x19,
	"ctrl+z":          0x1A,
	"esc":             0x1B,
	"ctrl+3":          0x1B,
	"ctrl+4":          0x1C,
	"ctrl+backslash":  0x1C,
	"ctrl+5":          0x1D,
	"ctrl+rsqbracket": 0x1D,
	"ctrl+6":          0x1E,
	"ctrl+7":          0x1F,
	"ctrl+slash":      0x1F,
	"ctrl+underscore": 0x1F,
	"space":           0x20,
	"backspace2":      0x7F,
	"ctrl+8":          0x7F,
}

func parse(c *Config) error {
	va, found := keylist[c.Monitor.CloseMonitorKey]
	if !found {
		return fmt.Errorf("can't find the close monitor key")
	}
	c.Monitor.CloseMonitorKeyUint16 = va

	va, found = keylist[c.Monitor.OpenCmdLineKey]
	if !found {
		return fmt.Errorf("can't find the open cmdline key")
	}
	c.Monitor.OpenCmdLineKeyUint16 = va

	return nil
}
