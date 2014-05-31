package main

import (
	"bufio"
	"net"
)

func (s *CacheServer) CommandReplace(conn net.Conn, reader *bufio.Reader, tokens []string, timestamp int64) {

	if len(tokens) != 5 {
		conn.Write([]byte("Error"))
		return
	}

	key := tokens[1]

	s.Store.RLock()

	if _, ok := s.Store.Data[key]; ok {
		s.Store.RUnlock()
		s.CommandSet(conn, reader, tokens, timestamp)
	} else {
		s.Store.RUnlock()
		conn.Write([]byte("NOT STORED\r\n"))
	}

}
