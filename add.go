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

	// TODO: There is a race condition here.
	// RLock, check if exists, Unlock, Send to
	// set, which grabs new lock. Can probably
	// sneak in a value between the Unlock => Set.
	//
	// Also don't want to cause a deadlock by holding
	// ReadLock.
	//
	// Maybe I should totally lock the data structure on
	// writes? Maybe can get around the suckage of this by
	// allocating a ring of diff maps to store values in.
	if _, ok := s.Store.Data[key]; !ok {
		s.Store.RUnlock()
		s.CommandSet(conn, reader, tokens)
	} else {
		s.Store.RUnlock()
		conn.Write([]byte("NOT STORED\r\n"))
	}

}
