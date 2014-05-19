package main

import (
  "net"
  "fmt"
  "log"
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

var setCommand  = []byte("s")[0]
var getCommand  = []byte("g")[0]
var quitCommand = []byte("q")[0]

func CacheServerRawHandler(conn net.Conn) {

  // TODO: Use a circular buffer or a static allocation
  // TODO: 1024 buffer is probably not enough
  buffer := make([]byte, 1024)
  defer conn.Close()

  // Loop and read, parsing for commands along the way
  for {
    count, err := conn.Read(buffer)

    if err != nil {
      log.Println("Error reading from client", err.Error())
      return
    }

    command := strings.TrimSpace(string(buffer))
    fmt.Println("got ", count, "bytes, ", command)

    if buffer[0] == quitCommand {
      return
    }

    if buffer[0] == setCommand {
      log.Println("set")
      conn.Write([]byte("STORED\r\n"))
    }

    if buffer[0] == getCommand {
      log.Println("get")
      conn.Write([]byte("VALUE key flags bytes\r\n"))
      conn.Write([]byte("Some data\r\n"))
    }

  }
}
