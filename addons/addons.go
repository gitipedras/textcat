package addons

import (
	"github.com/yuin/gopher-lua"
	"fmt"

	// internal
	"textcat/models"
)

var (
	L               *lua.LState
	commandRegistry = make(map[string]*lua.LFunction)
)


func AddonsInit() {
	L = lua.NewState()

	// Register Go function into Lua
	L.SetGlobal("coreRegisterCommand", L.NewFunction(coreRegisterCommand))
	L.SetGlobal("coreSendMessage", L.NewFunction(coreSendMessage))
	L.SetGlobal("coreKickUser", L.NewFunction(coreKickUser))

	// Load a Lua file
	if err := L.DoFile("commands.lua"); err != nil {
		panic(err)
	}
}

func Close() {
	if L != nil {
		L.Close()
	}
}

func HandleCommand(cmd string, data models.CommandData) {
	if fn, ok := commandRegistry[cmd]; ok {
		L.Push(fn)
		if err := L.PCall(0, lua.MultRet, nil); err != nil {
			fmt.Println("Lua command error:", err)
		}
	} else {
		fmt.Println("Unknown command:", cmd)
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
