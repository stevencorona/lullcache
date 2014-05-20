package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
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
	Exptime string
	Value   string
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
				conn.Write([]byte(item.Value))
				conn.Write([]byte("VALUE key flags bytes\r\n"))
				conn.Write([]byte("data blog\r\n"))
			}

			conn.Write([]byte("END\r\n"))
		}

		if command == "set" {

			if len(tokens) != 4 {
				conn.Write([]byte("Error"))
				return
			}

			key := tokens[1]
			flags := tokens[2]
			exptime := tokens[3]
			//bytes := tokens[4]

			cacheData[key] = CacheItem{exptime, "test", flags}

			conn.Write([]byte("STORED\r\n"))
		}

	}
}
