package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

type CacheItem struct {
	Exptime int64
	Value   []byte
	Flag    string
}

type CacheStore struct {
	*sync.RWMutex
	Data map[string]CacheItem
}

type CacheServer struct {
	Listener net.Listener
	Store    CacheStore
}

func NewCacheServer(address string) *CacheServer {
	store := CacheStore{&sync.RWMutex{}, make(map[string]CacheItem)}
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Error creating TCP socket", err.Error())
	}

	return &CacheServer{listener, store}
}

func (s *CacheServer) Start() {

	// Safe to close the listener after error checking
	defer s.Listener.Close()

	// Loop, accept, push work into a goroutine
	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			log.Println("Connection error from accept", err.Error())
		}

		// TODO: Use a pool of Goroutines
		go s.RawHandler(conn)
	}
}

func (s *CacheServer) RawHandler(conn net.Conn) {

	reader := bufio.NewReader(conn)
	defer conn.Close()

	// TODO extract this into an ASCIIProtocolHandler
	protocol := textproto.NewReader(reader)

	// Peek one byte and look for the magic byte that'll distinguish this as
	// a binary protocol connection
	peek, err := reader.Peek(1)

	if err != nil {
		log.Println("Error peeking at first byte", err.Error())
		return
	}

	if peek[0] == 0x80 {
		fmt.Println("binary")
	} else {
		fmt.Println("ascii")
	}

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

		if len(tokens) < 1 {
			log.Println("Command Error")
			return
		}

		// This should be dependent on the protocol instead of using magical
		// strings
		command := tokens[0]

		switch command {
		case "get", "gets":
			s.CommandGet(conn, tokens, timestamp)
		case "set":
			s.CommandSet(conn, reader, tokens, timestamp)
		case "delete":
			s.CommandDelete(conn, tokens)
		case "replace":
			s.CommandReplace(conn, reader, tokens, timestamp)
		case "add":
			s.CommandAdd(conn, reader, tokens, timestamp)
		case "touch":
			s.CommandTouch(conn, tokens, timestamp)
		case "quit":
			return
		}
	}
}
