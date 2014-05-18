package main

import (
  "net"
  "fmt"
  "log"
  "strings"
  )

func main() {
  ln, err := net.Listen("tcp", "127.0.0.1:11211")

  if err != nil {
    log.Fatal("Error creating TCP socket")
  }

  defer ln.Close()

  for {
    conn, err := ln.Accept()

    if err != nil {
      log.Println("Connection error on accept")
    }

    go ServerFunc(conn)

  }
}

func ServerFunc(conn net.Conn) {


  buf := make([]byte, 1024)
  defer conn.Close()

  for {
    n, err := conn.Read(buf)

    if err != nil {
      log.Println("error reading")
      return
    }

    command := strings.TrimSpace(string(buf))
    fmt.Println("got ", n, "bytes, ", command)

    quitCommand := []byte("q")[0]

    if (buf[0] == quitCommand) {
      return
    }

  }

}
