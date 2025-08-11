let logingui = document.getElementById("logingui")
let maingui = document.getElementById("maingui")
let userBox = document.getElementById("userBox")

maingui.style.display = 'none'
userBox.style.display = 'none'

function guiTransition() {
    if (logingui.style.display === 'block' || logingui.style.display === '') {
        // If login GUI is visible (or not explicitly hidden), show main GUI
        maingui.style.display = 'block';
        userBox.style.display = 'block';
        logingui.style.display = 'none';
    } else {
        // Otherwise, show login GUI and hide main GUI
        maingui.style.display = 'none';
        userBox.style.display = 'none';
        logingui.style.display = 'block';
    }
}


function preConnectRegister() {
    const registering = true
    preConnect(registering)
}

function preConnect(registering) {
    // URL's and page forms
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    //let input = document.getElementById("input").value;
    let sessionToken = document.getElementById("sessiontoken").value;

    let url = "ws://" + address + "/ws"
    prewebSocket = new WebSocket(url);


    prewebSocket.onopen = function() {
        console.log("Connected to PRE WebSocket server");

        if (( registering == undefined || registering == null )) {
            let payload = {
                rtype: "loginRequest",
                username: username,
                sessionToken: sessionToken,
                //clientID: "textcat_gui_client",
            };
            prewebSocket.send(JSON.stringify(payload));
        } else {
            let payload = {
                rtype: "register",
                username: username,
                sessionToken: sessionToken,
                //clientID: "textcat_gui_client",
            };
            prewebSocket.send(JSON.stringify(payload));
            console.log("Register success! Please login now.")
        }

    };

    prewebSocket.onmessage = function(event) {

        const msg = JSON.parse(event.data);
        console.log(msg)

        switch (msg.rtype) {
            case "invalidCredentials":
                alert("Your credentials are invalid!")
                break;

            case "goodCredentials":
                console.log("Login success")
                guiTransition()
                connect()

                break;

            case "rejected":
                alert("Login was rejected by the server")
                break;

            default:
                console.log("Unknown request type recieved from server: " + msg.rtype + " Possible reason: Outdated client or server!");
                break;
        }
    };

    prewebSocket.onclose = function() {
        console.log("PRE WebSocket connection closed.");
        //setTimeout(connect, 1000); // Reconnect after 1 second
    };

    prewebSocket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
}


function connect() {
    prewebSocket.close()

    // URL's and page forms
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    let input = document.getElementById("messageInput").value;
    let sessionToken = document.getElementById("sessiontoken").value;
    let usernamebox = document.getElementById("usernamebox")

    usernamebox.innerText = "Logged in as: " + username

    document.getElementById("messageInput").addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            sendMessage();
            this.value = "";
        }
    });

    let url = "ws://" + address + "/ws"
    webSocket = new WebSocket(url);


    webSocket.onopen = function() {
        console.log("Connected to WebSocket server");

        /*
        let payload = {
            rtype: "hello",
            username: username,
            sessionToken: sessionToken,
            message: "i am a client trying to connect!",
            //clientID: "textcat_gui_client",
        };
        webSocket.send(JSON.stringify(payload));
        */
        console.log("Websocket connection open")
    };


    webSocket.onmessage = function(event) {
        let messageDisplay = document.getElementById("messages");

        const msg = JSON.parse(event.data);
        console.log(msg)

        switch (msg.rtype) {
            case "invalidCredentials":
                alert("[Server] Your credentials are invalid!")
                break;

            case "goodCredentials":
                console.log("[Server] credentials OK")
                break;

            case "kicked":
                alert("[Server] Kicked by the server")
                break;

            case "sendMessage":
                messageDisplay.innerHTML += `<p>@${msg.username}: ${msg.message}</p>`;


            default:
                console.log("Unknown request type:" + msg.rtype + " Possible echoed request (BUG)");
                break;
        }

    };

    webSocket.onclose = function() {
        console.log("WebSocket connection closed, retrying...");
        //setTimeout(connect, 1000); // Reconnect after 1 second
    };

    webSocket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
}

function sendMessage() {
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    let input = document.getElementById("messageInput").value;
    let sessionToken = document.getElementById("sessiontoken").value;

    let payload = {
        rtype: "sendMessage",
        username: username,
        sessionToken: sessionToken,
        message: input,
        //clientID: "textcat_gui_client",
    };
    webSocket.send(JSON.stringify(payload));


    input.value = "";
}

function logout() {
    guiTransition();
    webSocket.close();
    alert("Logged out")
}
