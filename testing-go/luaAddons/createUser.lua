-- Register a new command
goRegisterCommand("/createUser", function(username)
    print("Creating user:", username)
    -- (you could call back into Go with another exposed function)
end)
