package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
)

func (s *CacheServer) CommandSet(conn net.Conn, reader *bufio.Reader, tokens []string, timestamp int64) {

	if len(tokens) != 5 {
		conn.Write([]byte("Error"))
		return
	}

	key := tokens[1]
	flags := tokens[2]
	exptime, _ := strconv.ParseInt(tokens[3], 10, 32)
	length, _ := strconv.ParseInt(tokens[4], 10, 32)

	// if the exptime is less than 30 days, it's probably just number of
	// seconds and not a timestamp, so add it to unix time.
	if exptime != 0 && exptime < 2592000 {
		exptime += timestamp
	}

	// Guard this
	bytes := make([]byte, length)
	io.ReadFull(reader, bytes)

	fmt.Println("got this:", string(bytes))

	s.Store.Lock()
	s.Store.Data[key] = CacheItem{exptime, bytes, flags}
	s.Store.Unlock()

	conn.Write([]byte("STORED\r\n"))
}
