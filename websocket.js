let logingui = document.getElementById("logingui");
let maingui = document.getElementById("maingui");
let userBox = document.getElementById("userBox");
let sidebar = document.getElementById("sidebar");
let chatInputBar = document.getElementById("chatInputBar");

var token
var chid

let savedTheme = localStorage.getItem("savedTheme")
document.body.setAttribute("data-theme", savedTheme)

const themeSelector = document.getElementById("themeForm");
const select = document.getElementById("themes");

themeSelector.addEventListener("submit", function (event) {
    event.preventDefault(); // stop page reload
    const theme = select.value; // get chosen theme
    document.body.setAttribute("data-theme", theme); // apply it
    localStorage.setItem("savedTheme", theme);
});

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
                rtype: "login",
                username: username,
                sessionToken: sessionToken,
                //clientID: "textcat_gui_client",
            };
            console.log(payload)
            prewebSocket.send(JSON.stringify(payload));
        } else {
            let payload = {
                rtype: "register",
                username: username,
                sessionToken: sessionToken,
                //clientID: "textcat_gui_client",
            };

            prewebSocket.send(JSON.stringify(payload));
        }

    };

    prewebSocket.onmessage = function(event) {

        const msg = JSON.parse(event.data);
        console.log(msg)

        switch (msg.Rtype) {
            case "loginStats":
                if (msg.Status == "ok") {
                    console.log("Login success, happy chatting ;)")
                    realClose = true

                    //showAlert("Loggin in...")

                    guiTransition()
                    token = msg.Value
                    connect(msg.Value)
                } else {
                    showAlert("Username or Password is incorrect")
                }

                break;

            case "registerStats":
                if (msg.Status == "ok") {
                    console.log("Register success, happy chatting ;)")
                    realClose = true

                    showAlert("Account created, please log in to activate the account")
                } 
                if (msg.Status == "isr") {
                     showAlert("Internal server error")
                } else {
                    showAlert("Username is taken")
                }

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


function connect(token) {
    prewebSocket.close()

    // URL's and page forms
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    let input = document.getElementById("messageInput").value;
    let sessionToken = document.getElementById("sessiontoken").value;
    let usernamebox = document.getElementById("usernamebox")

    usernamebox.innerHTML = `Logged in as: <span id="username-link" style="color: var(--accent-color); cursor: pointer;">${username}</span>`;

    // Add click handler
    document.getElementById("username-link").onclick = function() {
        showUser({
          username: username,
          token: token,
          description: "nothing here",
          dateCreated: "nothing here"
        });
    };


    document.getElementById("messageInput").addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            sendMessage(token);
            this.value = "";
        }
    });

    let url = "ws://" + address + "/ws"
    webSocket = new WebSocket(url);


    webSocket.onopen = function() {
        console.log("Connected to WebSocket server");
        console.log("Websocket connection open")
    };

    webSocket.onmessage = function(event) {
        let messageDisplay = document.getElementById("messages");
        console.log(event)

        const msg = JSON.parse(event.data);
        console.log("message is ", msg)
        console.log(msg.Rtype)

        switch (msg.Rtype) {
            case "invalidCredentials":
                showAlert("[Server] Your credentials are invalid!")
                break;

            case "goodCredentials":
                //console.log("[Server] credentials OK")
                break;

            case "kicked":
                showAlert("[Server] Kicked by the server");
                logout();
                break;

            case "invalidChannel":
                showAlert("[Server] Invalid Channel: " + msg.value)
                console.warn("Invalid channel: " + msg.value)
                break;

            case "unknownReq":
                console.warn("Recieved unknown request type from server (outdated client or server??)")
                break;

            /*case "isr":
                showAlert("[Server] Internal Server Error")
                logout()
                break;
            */

            case "alreadyConnected":
                showAlert("You are already connected to this channel!")
                break;

            case "invalidSession":
                showAlert("An invalid session was provided")
                logout()
                break;

            case "newMessage":
                console.log("New message: " + msg.Value)
                break;

            case "message":
                messageDisplay.innerHTML += `<p>@${msg.username}: ${msg.message}</p>`;
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

// Add this new function to your client.js file
function ConnectToChannel(channelName) {
    chid = channelName
    if (webSocket && webSocket.readyState === WebSocket.OPEN) {
        const payload = {
            Rtype: "connect",
            sessionToken: token,
            ChannelID: channelName
        };
        webSocket.send(JSON.stringify(payload));
        console.log(`Connection request sent for channel: ${channelName}`);
    } else {
        console.error("Cannot connect to channel. Connection not open or user info missing.");
    }
}

// Add this event listener to the bottom of your client.js file
document.addEventListener("DOMContentLoaded", () => {
    // ... (your existing code) ...

    const channelList = document.getElementById("channelList");

    channelList.addEventListener("click", (event) => {
        // Prevent the default link behavior (navigating to a new page)
        event.preventDefault();

        // Check if the clicked element is a channel link
        if (event.target.classList.contains("channel-link")) {
            // Get the channel name from the data-channel attribute
            const channelName = event.target.getAttribute("data-channel");
            ConnectToChannel(channelName);
        }
    });

    // ... (rest of your existing code) ...
});

function sendMessage(token) {
    let address = document.getElementById("server").value;
    let username = document.getElementById("username").value;
    let input = document.getElementById("messageInput").value;
    let sessionToken = document.getElementById("sessiontoken").value;
    let channelID = chid;

    let payload = {
        Rtype: "message",
        Username: username,
        SessionToken: token,
        Message: input,
        ChannelID: channelID,
        //clientID: "textcat_gui_client",
    };
    webSocket.send(JSON.stringify(payload));
    console.log(payload)

    input.value = "";
}

function logout() {
    guiTransition();
    webSocket.close();
    alert("Logged out")
}
