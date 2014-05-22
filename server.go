package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

func NewCacheServer(address string) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Error creating TCP socket", address, err.Error())
	}

	// Safe to close the listener after error checking
	defer listener.Close()

	// Loop, accept, push work into a goroutine
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Connection error from accept", err.Error())
		}

		// TODO: Use a pool of Goroutines
		go CacheServerRawHandler(conn)
	}
}

type CacheItem struct {
	Exptime int64
	Value   []byte
	Flag    string
}

var cacheData = make(map[string]CacheItem)

func CacheServerRawHandler(conn net.Conn) {

	reader := bufio.NewReader(conn)
	// TODO extract this into an ASCIIProtocolHandler
	protocol := textproto.NewReader(reader)

	defer conn.Close()

	// Loop and read, parsing for commands along the way
	for {

		line, err := protocol.ReadLine()

		timestamp := time.Now().Unix()
		log.Println(timestamp)

		if err != nil {
			log.Println("Error reading from client", err.Error())
			return
		}

		fmt.Println("got line: ", line)

		tokens := strings.Split(line, " ")

		// TODO: Guard
		command := tokens[0]

		if command == "quit" {
			return
		}

		if command == "get" {
			// TODO: Guard this
			key := tokens[1]

			if item, ok := cacheData[key]; ok {

				if timestamp > item.Exptime {
					log.Println("expiring key:", key)
					delete(cacheData, key)
				} else {
					out := fmt.Sprintf("VALUE %s %s %d\r\n%s\r\n", key, item.Flag, item.Exptime, item.Value)
					conn.Write([]byte(out))
				}
			}

			conn.Write([]byte("END\r\n"))
		}

		if command == "set" {

			if len(tokens) != 5 {
				conn.Write([]byte("Error"))
				return
			}

			key := tokens[1]
			flags := tokens[2]
			exptime, _ := strconv.ParseInt(tokens[3], 10, 32)
			length, _ := strconv.ParseInt(tokens[4], 10, 32)

			if exptime < 2592000 {
				exptime += timestamp
			}

			// Guard this
			bytes := make([]byte, length)
			io.ReadFull(reader, bytes)

			fmt.Println("got this:", string(bytes))

			cacheData[key] = CacheItem{exptime, bytes, flags}

			conn.Write([]byte("STORED\r\n"))
		}

	}
}
