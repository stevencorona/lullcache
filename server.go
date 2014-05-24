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
		go s.CacheServerRawHandler(conn)
	}
}

func (s *CacheServer) CacheServerRawHandler(conn net.Conn) {

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

		if len(tokens) < 1 {
			log.Println("Command Error")
			return
		}

		command := tokens[0]

		if command == "quit" {
			return
		}

		if command == "get" || command == "gets" {
			s.CommandGet(conn, tokens, timestamp)
		}

		if command == "set" {
			s.CommandSet(conn, reader, tokens, timestamp)
		}

		if command == "replace" {

		}

		if command == "add" {

		}

		if command == "cas" {

		}

		if command == "prepend" {

		}

		if command == "append" {

		}

	}
}

func (s *CacheServer) CommandGet(conn net.Conn, tokens []string, timestamp int64) {
	for _, key := range tokens[1:] {

		s.Store.RLock()
		if item, ok := s.Store.Data[key]; ok {

			s.Store.RUnlock()

			if timestamp > item.Exptime {
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
