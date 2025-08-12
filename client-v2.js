let logingui = document.getElementById("logingui");
let maingui = document.getElementById("maingui");
let userBox = document.getElementById("userBox");
let sidebar = document.getElementById("sidebar");
let chatInputBar = document.getElementById("chatInputBar");

// Initially hide main UI components except login
maingui.style.display = 'none';
userBox.style.display = 'none';
sidebar.style.display = 'none';
chatInputBar.style.display = 'none';

function guiTransition() {
    if (logingui.style.display === 'block' || logingui.style.display === '') {
        // Switch from login to main UI
        maingui.style.display = 'flex';      // Use flex because layout uses flexbox
        userBox.style.display = 'flex';
        sidebar.style.display = 'flex';
        chatInputBar.style.display = 'flex';
        logingui.style.display = 'none';
    } else {
        // Switch from main UI back to login
        maingui.style.display = 'none';
        userBox.style.display = 'none';
        sidebar.style.display = 'none';
        chatInputBar.style.display = 'none';
        logingui.style.display = 'block';
    }
}



function preConnectRegister() {
    const registering = true
    preConnect(registering)
}

let realClose = false;

function preConnect(registering, realClose) {
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

            case "registerOk":
                alert("[Server] Register success!")
                break;

            case "internalServerErr":
                alert("[Server] Internal Server Error")
                break;

            case "alreadyExists":
                alert("[Server] Username is already taken")
                break;

            case "goodCredentials":
                console.log("Login success")
                realClose = true
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
        let Premsg = document.getElementById("Premsg")

        console.log("PRE WebSocket connection closed.");
        Premsg.innerText = "Server unreachable... Retrying..."

        if (( realClose == null && realClose == undefined)) {

            console.log("Reconnecting...")
            setTimeout(preConnect, 1000); // Reconnect after 1 second

        } else {
            console.log("RealClose is disabled.")
        }
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
                alert("[Server] Kicked by the server");
                logout();
                break;

            case "unknownReq":
                console.warn("Recieved unknown request type from server (outdated client or server??)")

            case "internalServerErr":
                alert("[Server] Internal Server Error")

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
