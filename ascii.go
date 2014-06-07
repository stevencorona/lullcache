package main

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
