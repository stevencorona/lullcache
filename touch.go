package main

import (
	"net"
	"strconv"
	"time"
)

func (s *CacheServer) CommandTouch(conn net.Conn, tokens []string) {

	if len(tokens) != 3 {
		conn.Write(ERROR)
		return
	}

	key := tokens[1]
	exptime, _ := strconv.ParseInt(tokens[2], 10, 32)

	// if the exptime is less than 30 days, it's probably just number of
	// seconds and not a timestamp, so add it to unix time.
	if exptime != 0 && exptime < 2592000 {
		exptime += time.Now().Unix()
	}

	if item, ok := s.Store.Data[key]; ok {

		s.Store.Lock()
		item.Exptime += exptime
		s.Store.Data[key] = item
		s.Store.Unlock()
		conn.Write(TOUCHED)
	} else {
		conn.Write(NOT_FOUND)
	}
}
