package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

var commandRegistry = make(map[string]*lua.LFunction)

func main() {
	L := lua.NewState()
	defer L.Close()

	// Register Go function into Lua
	L.SetGlobal("goRegisterCommand", L.NewFunction(goRegisterCommand))

	// Load a Lua file
	if err := L.DoFile("commands.lua"); err != nil {
		panic(err)
	}

	// Simulate someone typing "/hi"
	handleCommand(L, "hi")
}

func goRegisterCommand(L *lua.LState) int {
	cmdName := L.CheckString(1)     // first arg: command name
	fn := L.CheckFunction(2)        // second arg: Lua function
	commandRegistry[cmdName] = fn   // store it in registry
	return 0
}
