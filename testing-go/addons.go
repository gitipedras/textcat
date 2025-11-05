package main

import (
	//"fmt"
	"github.com/yuin/gopher-lua"
)

var commandRegistry = make(map[string]*lua.LFunction)

func main() {
	L := lua.NewState()
	defer L.Close()

	// Register Go function into Lua
	L.SetGlobal("coreRegisterCommand", L.NewFunction(coreRegisterCommand))
	L.SetGlobal("coreSendMessage", L.NewFunction(coreSendMessage))
	L.SetGlobal("coreKickUser", L.NewFunction(coreKickUser))

	// Load a Lua file
	if err := L.DoFile("commands.lua"); err != nil {
		panic(err)
	}

	// Simulate someone typing "/hi"
	handleCommand(L, "hi")
}

func handleCommand(L *lua.LState, cmd string) {
	// Look up the command in the registry
	if fn, ok := commandRegistry[cmd]; ok {
		// Push the function onto the Lua stack
		L.Push(fn)

		// Call the Lua function (no arguments, 0 return values)
		if err := L.PCall(0, lua.MultRet, nil); err != nil {
			panic(err)
		}
	} else {
		println("Unknown command:", cmd)
	}
}

func coreRegisterCommand(L *lua.LState) int {
	cmdName := L.CheckString(1)     // first arg: command name
	fn := L.CheckFunction(2)        // second arg: Lua function
	commandRegistry[cmdName] = fn   // store it in registry
	return 0
}

func coreSendMessage(L *lua.LState) int {
	message := L.CheckString(1)
	channelid := L.CheckString(2)

	fmt.Printf(message)
	fmt.Printf(channelid)

	return 0
}

func coreKickUser(L *lua.LState) int {
	username := L.CheckString(1)

	fmt.Printf(username)

	return 0
}