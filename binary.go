package main

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

var AsciiCommands = map[string]int{
	"get":  Get,
	"gets": Get,
	"set":  Set,
}
