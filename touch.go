package main

import (
	"net"
	"strconv"
)

func (s *CacheServer) CommandTouch(conn net.Conn, tokens []string, timestamp int64) {

	key := tokens[1]
	exptime, _ := strconv.ParseInt(tokens[2], 10, 32)

	// if the exptime is less than 30 days, it's probably just number of
	// seconds and not a timestamp, so add it to unix time.
	if exptime != 0 && exptime < 2592000 {
		exptime += timestamp
	}

	if item, ok := s.Store.Data[key]; ok {

		s.Store.Lock()
		item.Exptime += exptime
		s.Store.Data[key] = item
		s.Store.Unlock()
		conn.Write([]byte("TOUCHED\r\n"))
	} else {
		conn.Write([]byte("NOT FOUND\r\n"))
	}
}
