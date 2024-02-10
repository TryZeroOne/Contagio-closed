package monitor

import (
	"contagio/contagio/config"
	"contagio/contagio/database/sqlite"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintTable(c *config.Config) {

	var data [][]string
	var traffic string
	fmt.Println(config.ParseColors("Press {blue}" + c.Monitor.OpenCmdLineKey + "{reset} to open the command line"))

	isStarted, status_message := isStarted(c)
	if isStarted {

		pid := sqlite.GetPid()
		if pid == "" {
			return
		}

		sessions := sqlite.GetSessions()
		users := sqlite.GetCount()
		bots := sqlite.GetBots()
		inc, out := sqlite.GetStats()

		data = [][]string{
			[]string{status_message, pid, strconv.Itoa(users), sessions, bots},
		}

		i, _ := strconv.Atoi(inc)
		o, _ := strconv.Atoi(out)
		traffic = fmt.Sprintf("Incoming: %s\nOutgoing: %s\n", parseTraffic(float64(i)/1024.0), parseTraffic(float64(o)/1024.0))

	} else {
		data = [][]string{
			// []string{}
			[]string{status_message, "NULL", "NULL", "NULL", "NULL"},
		}
		traffic = fmt.Sprintf("Incoming: %s\nOutgoing: %s\n", "NULL", "NULL")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CNC STATUS", "CNC PID", "USERS", "SESSIONS", "BOTS"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetRowSeparator("─")
	table.SetCenterSeparator("─")
	table.SetColumnSeparator("│")
	table.SetCenterSeparator("┼")
	table.AppendBulk(data)
	table.Render()

	fmt.Println("MEMORY:", parseMemory(), "\nCPU:", parseCpu())
	fmt.Println(traffic)

}
