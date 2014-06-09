package main

import (
	"net"
)

func (s *CacheServer) CommandDelete(conn net.Conn, tokens []string) {

	if len(tokens) != 2 {
		conn.Write(ERROR)
		return
	}

	key := tokens[1]

	s.Store.RLock()

	if _, ok := s.Store.Data[key]; ok {
		s.Store.RUnlock()

		s.Store.Lock()
		delete(s.Store.Data, key)
		s.Store.Unlock()

		conn.Write(DELETED)
	} else {
		s.Store.RUnlock()
		conn.Write(NOT_FOUND)
	}
}
