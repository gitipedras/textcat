let maingui = document.getElementById("maingui")
maingui.style.display = 'none'

function login() {

    let maingui = document.getElementById("maingui")
    maingui.style.display = 'block'

    let logingui = document.getElementById("logingui")
    logingui.style.display = 'none'

    connect();

}

let prews;
let prewsClosedIntentionally = false;

function preConnect() {
    let statusmsg = ""

    let serveraddress = document.getElementById("server").value
    let fulladdress = "ws://" + serveraddress + "/ws"

    prews = new WebSocket(fulladdress);

    prews.onopen = function() {
        console.log("Connected to WebSocket server... logging in...");
        document.getElementById("Premsg").innerText = statusmsg + "Connected to server!"

        sendLoginRequest()

        maingui.style.display = 'block'
        logingui.style.display = 'none'
    };

    prews.onmessage = function(event) {
        console.log(event.data)
    };

    prews.onclose = function() {
        console.log("WebSocket login connection closed");

        if (!prewsClosedIntentionally) {
            console.log("Retrying...");
            document.getElementById("Premsg").innerText = statusmsg + "Attempting to connect..."
            setTimeout(preConnect, 1000); // Only retry if not intentional
        }
    };

    prews.onerror = function(error) {
        console.error("WebSocket error:", error);
    };

}

function sendLoginRequest() {
    let username = document.getElementById("username").value;
    let sessiontoken = document.getElementById("sessiontoken").value;

    let message = "loginRequest" + "," + username + "," + sessiontoken

    prews.send(message);
    prewsClosedIntentionally = true; // <- prevent reconnect
    prews.close();

    connect()

}


let ws;
let wsClosedIntentionally = false;

function connect() {
    let serveraddress = document.getElementById("server").value
    let fulladdress = "ws://" + serveraddress + "/ws"
    let statusmsg = "Status: "

    ws = new WebSocket(fulladdress);

    ws.onopen = function() {
        console.log("Connected to WebSocket server");
        document.getElementById("msg").innerText = statusmsg + "Connected to server!"
    };

    ws.onmessage = function(event) {
        let messageDisplay = document.getElementById("messages");
        messageDisplay.innerHTML += `<p>${event.data}</p>`;
    };

    ws.onclose = function() {
        console.log("WebSocket login connection closed");

        if (!wsClosedIntentionally) {
            console.log("Retrying...");
            document.getElementById("Premsg").innerText = statusmsg + "Attempting to connect..."
            setTimeout(preConnect, 1000); // Only retry if not intentional
        }
    };

    ws.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
}

function logout() {

    ws.close()
    wsClosedIntentionally = true
}

function sendMessage() {
    let input = document.getElementById("messageInput");
    let username = document.getElementById("username").value;
    let sessiontoken = document.getElementById("sessiontoken").value;

    let message = "sendMessage" + "," + input.value + "," + sessiontoken

    ws.send(message);
    input.value = "";
}
