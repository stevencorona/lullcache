package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func (s *CacheServer) CommandGet(conn net.Conn, tokens []string) {
	for _, key := range tokens[1:] {

		s.Store.RLock()
		if item, ok := s.Store.Data[key]; ok {

			s.Store.RUnlock()

			// Get the current timestamp
			timestamp := time.Now().Unix()

			if timestamp != 0 && timestamp > item.Exptime {
				log.Println("expiring key:", key)

				s.Store.Lock()
				delete(s.Store.Data, key)
				s.Store.Unlock()
			} else {
				out := fmt.Sprintf("VALUE %s %s %d\r\n%s\r\n", key, item.Flag, item.Exptime, item.Value)
				conn.Write([]byte(out))
			}
		}
	}

	conn.Write([]byte("END\r\n"))
}
