var userToken
let currentChannel = "main";
let msgInput = document.getElementById("messageInput")
let username = document.getElementById("username").value
alreadyRan = false

function wsConnect(action, address, password) {
// action is if ur signin in or registering
// msg is the object with the input field
// address is the object with the server address

    let url = "ws://" + address + "/ws"
    webSocket = new WebSocket(url);

    let sidebar = document.getElementById("sidebar")
    function inputInit() {
        console.log("input init")
        document.querySelectorAll(".channel-link").forEach(link => {
            link.addEventListener("click", e => {
                e.preventDefault(); // stop the "#" jump
                clearChat()

                if (alreadyRan == true) {
                    disconnectChannel(currentChannel);
                }

                const channelName = link.getAttribute("data-channel");
                connectChannel(channelName);
            });
        });
        msgInput.addEventListener("keydown", function(e) {
            if (e.key === "Enter") {  // check if Enter was pressed
                e.preventDefault();   // optional: prevent default form submission
                writeMessage();        

                messageInput.value = ""; // clear input after sending
            }
        });
    }


    webSocket.onopen = function() {
        console.log("Connected to WebSocket server");
        console.log("Websocket connection open")

        if (action == "login") {
                let payload = {
                    Rtype: "login",
                    Username: username,
                    SessionToken: password,
                };
            console.log("Sent login request to server at", address)
            webSocket.send(JSON.stringify(payload));

        } else {
            let payload = {
                Rtype: "register",
                Username: username,
                SessionToken: password,
            };
            console.log("Sent register request to server at", address)
            webSocket.send(JSON.stringify(payload));
        }
    };
    

    webSocket.onmessage = function(event) {

        const msg = JSON.parse(event.data);
        console.log("message is ", msg)

        switch (msg.Rtype) {
            /* --- Loggin In and Registering */
            case "loginStats":
                if (msg.Status == "ok") {
                    userToken = msg.Value

                    setDetails(msg.ServerName, msg.ServerDesc)
                    guiTransition()
                    inputInit()
                } else if (msg.Status == "invalidInput") {
                    showAlert("Invalid username: username must only contain normal letters (capital included), numbers and underscores. Dashes ('-') are not supported.")

                } else {
                    showAlert("Invalid username or password")
                }
            break;

            case "registerStats":
                if (msg.Status == "ok") {
                    showAlert("Register Success, please login now.")
                } else {
                    showAlert("Username is taken")
                }
            break;

            case "invalidInput":
                showAlert("Invalid Input: Messages cannot be empty or longer than 70 characters")
                break;
                
            case "kicked":
                showAlert("Connection force-closed by the server");
                logout();
                break;

            case "rejected":
                showAlert("Login was rejected by the server")
                break;

            /* --- Client stuff --- */
            case "disconnectStats":
                if (msg.Status == "NoChannelFound") {
                    showAlert("Channel doesn't exist: " + msg.Value)
                } if (msg.Status == "notConnected" ) {
                    showAlert("Not connected to channel: " + msg.Value)
                } if (msg.Status == "ok" ) {
                    console.log("disconnect ok")
                } else {
                    showAlert("Failed to disconnect, unknown reason. Possible internal server error")
                }
                break;

            case "invalidChannel":
                showAlert("Invalid Channel")
                console.warn("Invalid channel")
                break;

            case "unknownReq":
                console.warn("Outdated client or server, server sent unknownReq")
                break;

            /*case "isr":
                showAlert("[Server] Internal Server Error")
                logout()
                break;
            */

            case "invalidToken":
                showAlert("An invalid token was provided to the server -- please login again.")
                break;
                
            case "alreadyConnected":
                showAlert("You are already connected to a channel")
                break;

            case "connectOk":
                console.log("connection to channel ok")
                break;

            case "failed":
                showAlert("Failed to send message, reason: ", msg.Status)
                break;


            case "invalidSession":
                showAlert("An invalid session was provided")
                logout()
                break;

            case "NewMessage":
                console.log("New message: " + msg.Value, " Sent by: " + msg.Username)
                messageDisplay(msg.Username, msg.Value)
                break;

            default:
                console.log("Unknown request type from server: " + msg.Rtype);
                break;
        }

    };

    webSocket.onclose = function() {
        showAlert("Connection closed.");
    };

    webSocket.onerror = function(error) {
        showAlert("Websocket Error: ", error)
        console.error("WebSocket error:", error);
    };
}

function logout() {
    showAlert("Log out")
    guiTransition()
    webSocket.close();
}

function disconnectChannel() {
    console.log("current channel is " + currentChannel)
    msgValue = msgInput.value
    let payload = {
                Rtype: "disconnect",
                SessionToken: userToken,
                ChannelID: currentChannel,
                Username: username,
    };
    console.log("Disconnecting: ", payload)
    webSocket.send(JSON.stringify(payload));
}

function connectChannel(channel) {
    alreadyRan = true
    msgValue = msgInput.value
    let payload = {
                Rtype: "connect",
                SessionToken: userToken,
                ChannelID: channel,
                Username: username,
    };
    currentChannel = channel
    console.log("Connected to channel: ", payload)
    webSocket.send(JSON.stringify(payload));
}

function writeMessage() {
    msgValue = msgInput.value
    let payload = {
                Rtype: "message",
                SessionToken: userToken,
                ChannelID: currentChannel,
                Message: msgValue,
                Username: username,
    };
    console.log("Sent message: ", payload)
    webSocket.send(JSON.stringify(payload));
}

function escapeHTML(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML; // safely escaped
}

function formatMessage(text) {
    // Escape first (prevents <script> etc.)
    let safe = escapeHTML(text);

    // Markdown replacements
    safe = safe
        .replace(/^#### (.*$)/gim, "<h4>$1</h4>")  // #### heading4
        .replace(/^### (.*$)/gim, "<h3>$1</h3>")  // ### heading3
        .replace(/\*\*(.*?)\*\*/g, "<b>$1</b>")   // **bold**
        .replace(/\*(.*?)\*/g, "<i>$1</i>")      // *italic*
        .replace(/#/g, "<br>");  // single # becomes a newline

    return safe;
}

function messageDisplay(username, message) {
    const messagesDiv = document.getElementById("messages");

    const wrapper = document.createElement("div");

    const userEl = document.createElement("b");
    userEl.textContent = username + ": ";

    const msgEl = document.createElement("span");
    msgEl.innerHTML = formatMessage(message); // safe + markdown

    wrapper.appendChild(userEl);
    wrapper.appendChild(msgEl);

    messagesDiv.appendChild(wrapper);
}



function setDetails(name, desc) {
    const sidebar = document.getElementById("sidebar");

    // create footer container
    const footer = document.createElement("div");
    footer.id = "sidebar-footer"; // optional id for styling

    // create name paragraph
    const nameP = document.createElement("p");
    nameP.textContent = name;
    footer.appendChild(nameP);

    // create description paragraph
    const descP = document.createElement("p");
    descP.textContent = desc;
    footer.appendChild(descP);

    // append footer to sidebar
    sidebar.appendChild(footer);
}

function clearChat() {
    document.getElementById("messages").innerHTML = "";
}
