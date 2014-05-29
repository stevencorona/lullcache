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
		go s.RawHandler(conn)
	}
}

func (s *CacheServer) RawHandler(conn net.Conn) {

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

		if command == "delete" {
			s.CommandDelete(conn, tokens)
		}

		if command == "replace" {
			s.CommandReplace(conn, reader, tokens, timestamp)
		}

		if command == "add" {
			s.CommandAdd(conn, reader, tokens, timestamp)
		}

		if command == "cas" {

		}

		if command == "prepend" {

		}

		if command == "append" {

		}

		if command == "touch" {
			s.CommandTouch(conn, tokens, timestamp)
		}

	}
}

func (s *CacheServer) CommandGet(conn net.Conn, tokens []string, timestamp int64) {
	for _, key := range tokens[1:] {

		s.Store.RLock()
		if item, ok := s.Store.Data[key]; ok {

			s.Store.RUnlock()

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

func (s *CacheServer) CommandTouch(conn net.Conn, tokens []string, timestamp int64) {

	key := tokens[1]
	exptime, _ := strconv.ParseInt(tokens[2], 10, 32)

	// if the exptime is less than 30 days, it's probably just number of
	// seconds and not a timestamp, so add it to unix time.
	if exptime != 0 && exptime < 2592000 {
		exptime += timestamp
	}

	if item, ok := s.Store.Data[key]; ok {

		s.Store.Lock()
		item.Exptime += exptime
		s.Store.Data[key] = item
		s.Store.Unlock()
		conn.Write([]byte("TOUCHED\r\n"))
	} else {
		conn.Write([]byte("NOT FOUND\r\n"))
	}
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

func (s *CacheServer) CommandAdd(conn net.Conn, reader *bufio.Reader, tokens []string, timestamp int64) {

	if len(tokens) != 5 {
		conn.Write([]byte("Error"))
		return
	}

	key := tokens[1]

	s.Store.RLock()

	if _, ok := s.Store.Data[key]; !ok {
		s.Store.RUnlock()
		s.CommandSet(conn, reader, tokens, timestamp)
	} else {
		s.Store.RUnlock()
		conn.Write([]byte("NOT STORED\r\n"))
	}

}

func (s *CacheServer) CommandReplace(conn net.Conn, reader *bufio.Reader, tokens []string, timestamp int64) {

	if len(tokens) != 5 {
		conn.Write([]byte("Error"))
		return
	}

	key := tokens[1]

	s.Store.RLock()

	if _, ok := s.Store.Data[key]; ok {
		s.Store.RUnlock()
		s.CommandSet(conn, reader, tokens, timestamp)
	} else {
		s.Store.RUnlock()
		conn.Write([]byte("NOT STORED\r\n"))
	}

}
