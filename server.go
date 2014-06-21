package main

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type CacheItem struct {
	Exptime int64
	Value   []byte
	Flag    string
}

type CacheStore struct {
	sync.RWMutex
	Data map[string]CacheItem
}

type CacheServer struct {
	Listener net.Listener
	Store    CacheStore
}

type Command struct {
	Magic    byte
	Opcode   int
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

func NewCacheServer(address string) *CacheServer {
	store := CacheStore{sync.RWMutex{}, make(map[string]CacheItem)}
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

	protocol := new(AsciiProtocol)

	if peek[0] == BINARY_MAGIC {
		log.Println("Protocol is Binary")
	} else {
		log.Println("Protocol is ASCII")
		// ascii
	}

	// Loop and read, parsing for commands along the way
	for {

		command, tokens := protocol.ReadCommand(reader)

		switch command.Opcode {
		case Get:
			s.CommandGet(conn, tokens)
		case Set:
			s.CommandSet(conn, reader, tokens)
		case Delete:
			s.CommandDelete(conn, tokens)
		case Replace:
			s.CommandReplace(conn, reader, tokens)
		case Add:
			s.CommandAdd(conn, reader, tokens)
		case Touch:
			s.CommandTouch(conn, tokens)
		case Quit:
			return
		}
	}
}
