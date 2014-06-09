package main

var STORED = []byte("STORED\r\n")
var DELETED = []byte("DELETED\r\n")
var NOT_FOUND = []byte("NOT FOUND\r\n")
var VALUE = "VALUE %s %s %d\r\n%s\r\n"
var TOUCHED = []byte("TOUCHED\r\n")
var END = []byte("END\r\n")

var AsciiCommands = map[string]int{
	"get":     Get,
	"gets":    Get,
	"set":     Set,
	"add":     Add,
	"touch":   Get,
	"delete":  Delete,
	"replace": Replace,
	"quit":    Quit,
}

func AsciiHandler() {

}
