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

type Command struct {
	Magic    byte
	Opcode   byte
	Length   byte
	Extra    byte
	Type     byte
	Reserved byte
	Body     []byte
	Opaque   []byte
	Cas      []byte
	Extras   []byte
	Key      []byte
	Value    []byte
}

// opcodes
var Get = 0x00
var Set = 0x01
var Add = 0x02
var Replace = 0x03
var Delete = 0x04
var Increment = 0x05
var Decrement = 0x06
var Quit = 0x07
var Flush = 0x08
var GetQ = 0x09
var Noop = 0x0A
var Version = 0x0B
var GetK = 0x0C
var GetKq = 0x0D
var Append = 0x0E
var Prepend = 0x0F
var Stat = 0x10
var SetQ = 0x11
var AddQ = 0x12
var ReplaceQ = 0x13
var DeleteQ = 0x14
var IncrementQ = 0x15
var DecrementQ = 0x16
var QuitQ = 0x17
var FlushQ = 0x18
var AppendQ = 0x19
var PrenendQ = 0x1A

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

	// Peek one byte and look for the magic byte that'll distinguish this as
	// a binary protocol connection
	peek, err := reader.Peek(1)

	if err != nil {
		log.Println("Error peeking at first byte", err.Error())
		return
	}

	protocol := textproto.NewReader(reader)

	if peek[0] == 0x80 {
		// binary
	} else {
		// ascii
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
			s.CommandGet(conn, tokens)
		case "set":
			s.CommandSet(conn, reader, tokens)
		case "delete":
			s.CommandDelete(conn, tokens)
		case "replace":
			s.CommandReplace(conn, reader, tokens)
		case "add":
			s.CommandAdd(conn, reader, tokens)
		case "touch":
			s.CommandTouch(conn, tokens)
		case "quit":
			return
		}
	}
}
