package main

import (
	"bufio"
	"net"
)

func (s *CacheServer) CommandAdd(conn net.Conn, reader *bufio.Reader, tokens []string) {

	if len(tokens) != 5 {
		conn.Write([]byte("Error"))
		return
	}

	key := tokens[1]

	s.Store.RLock()

	if _, ok := s.Store.Data[key]; !ok {
		s.Store.RUnlock()
		s.CommandSet(conn, reader, tokens)
	} else {
		s.Store.RUnlock()
		conn.Write([]byte("NOT STORED\r\n"))
	}

}
