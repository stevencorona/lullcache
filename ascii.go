package main

import (
	"bufio"
	"log"
	"strings"
)

var STORED = []byte("STORED\r\n")
var DELETED = []byte("DELETED\r\n")
var NOT_FOUND = []byte("NOT FOUND\r\n")
var VALUE = "VALUE %s %s %d\r\n%s\r\n"
var TOUCHED = []byte("TOUCHED\r\n")
var END = []byte("END\r\n")
var ERROR = []byte("ERROR\r\n")
var NOT_STORED = []byte("NOT STORED\r\n")

var AsciiCommands = map[string]int{
	"get":     Get,
	"gets":    Get,
	"set":     Set,
	"add":     Add,
	"touch":   Touch,
	"delete":  Delete,
	"replace": Replace,
	"quit":    Quit,
}

type AsciiProtocol struct {
}

func (ascii *AsciiProtocol) ReadCommand(reader *bufio.Reader) (*Command, []string) {

	command := new(Command)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Println("Error reading from client", err.Error())
		return command, nil
	}

	tokens := strings.Split(line, " ")

	if _, ok := AsciiCommands[tokens[0]]; ok {
		command.Opcode = AsciiCommands[tokens[0]]
	}

	return command, tokens

}
