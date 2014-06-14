package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func (s *CacheServer) CommandSet(conn net.Conn, reader *bufio.Reader, tokens []string) {

	if len(tokens) != 5 {
		conn.Write(ERROR)
		return
	}

	key := tokens[1]
	flags := tokens[2]
	exptime, expErr := strconv.ParseInt(tokens[3], 10, 32)
	length, lenErr := strconv.ParseInt(tokens[4], 10, 32)

	if expErr != nil {
		exptime = 0
	}

	if lenErr != nil {
		conn.Write(ERROR)
		return
	}

	// if the exptime is less than 30 days, it's probably just number of
	// seconds and not a timestamp, so add it to unix time.
	if exptime != 0 && exptime < 2592000 {
		exptime += time.Now().Unix()
	}

	// We're trusting the client to not lie to us about the size of the bytes.
	// This needs to be guarded from the client passing in more bytes than they
	// say.

	bytes := make([]byte, length)
	io.ReadFull(reader, bytes)

	log.Println("Server Received:", string(bytes))

	s.Store.Lock()
	s.Store.Data[key] = CacheItem{exptime, bytes, flags}
	s.Store.Unlock()

	conn.Write(STORED)
}
