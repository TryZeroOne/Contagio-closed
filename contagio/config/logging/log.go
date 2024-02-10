package logging

import (
	"contagio/contagio/config"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LogLevelInfo = 1 << iota
	LogLevelDiscard
	LogLevelWarning
	LogLevelError
)

var logLevel = LogLevelInfo

func SetLogLevels(levels string) {
	levels = strings.ReplaceAll(levels, "ERROR", strconv.Itoa(LogLevelError))
	levels = strings.ReplaceAll(levels, "WARNING", strconv.Itoa(LogLevelWarning))
	levels = strings.ReplaceAll(levels, "INFO", strconv.Itoa(LogLevelInfo))
	levels = strings.ReplaceAll(levels, "DISCARD", strconv.Itoa(LogLevelDiscard))

	res, err := calculate(levels)
	if err != nil {
		fmt.Println("Can't set log level: " + err.Error())
		os.Exit(0)
	}

	logLevel = res
}

func PrintInfo(message string) {
	if logLevel&LogLevelInfo != 0 {
		fmt.Println(parse(message, config.Config_.Logging.InfoMessage))
	}
}

func PrintWarning(message string) {
	if logLevel&LogLevelWarning != 0 {
		fmt.Println(parse(message, config.Config_.Logging.WarningMessage))
	}
}

func PrintError(message string) {
	if logLevel&LogLevelError != 0 {
		fmt.Println(parse(message, config.Config_.Logging.ErrorMessage))
	}
}

func parse(message string, template string) (result string) {

	result = strings.ReplaceAll(template, "{date}", time.Now().Format("15:04:05"))
	result = strings.ReplaceAll(result, "{message}", message)
	result = config.ParseColors(result)

	return result
}

func calculate(expression string) (int, error) {
	tokens := strings.Split(expression, "|")

	if len(tokens) == 0 {
		return 0, fmt.Errorf("null expression")
	}

	result, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(tokens); i++ {
		operand, err := strconv.Atoi(tokens[i])
		if err != nil {
			return 0, err
		}

		result |= operand
	}

	return result, nil
}
