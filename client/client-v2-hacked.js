let logingui = document.getElementById("logingui")
let maingui = document.getElementById("maingui")

maingui.style.display = 'none'

function guiTransition() {
    maingui.style.display = 'block'
    logingui.style.display = 'none'
}

function preConnect() {
    // URL's and page forms
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    //let input = document.getElementById("input").value;
    let sessionToken = document.getElementById("sessiontoken").value;

    let url = "ws://" + address + "/ws"
    prewebSocket = new WebSocket(url);


    prewebSocket.onopen = function() {
        console.log("Connected to PRE WebSocket server");

        let payload = {
            rtype: "loginRequest",
            username: username,
            sessionToken: sessionToken,
            //clientID: "textcat_gui_client",
        };
        prewebSocket.send(JSON.stringify(payload));

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
        messageDisplay.innerHTML += `<p>${event.data}</p>`;

        const msg = JSON.parse(event.data);
        console.log(msg)

        switch (msg.rtype) {
            case "invalidCredentials":
                alert("[Server] Invalid Credentials")
                console.log("[Server] Invalid Credentials")
                break;

            case "goodCredentials":
                console.log("[Server] Credentials OK")
                break;

            case "kicked":
                alert("[Server] Kicked by the server")
                break;

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
        username: "testuser",
        sessionToken: "abcd",
        message: input,
        //clientID: "textcat_gui_client",
    };
    webSocket.send(JSON.stringify(payload));


    input.value = "";
}

function logout() {

}
