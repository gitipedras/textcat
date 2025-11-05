var userToken
let currentChannel = "main";
let msgInput = document.getElementById("messageInput")
alreadyRan = false

// A named function to handle the channel link click event
const channelLinkClickHandler = e => {
    e.preventDefault(); // stop the "#" jump
    clearChat();

    // The 'this' keyword inside an arrow function needs to be handled
    // We get the channel name directly from the event target
    const link = e.currentTarget; 

    if (alreadyRan === true) {
        disconnectChannel(currentChannel);
    }

    const channelName = link.getAttribute("data-channel");
    connectChannel(channelName);

    // make clicked link bold
    document.querySelectorAll(".channel-link").forEach(l => l.classList.remove("active-channel"));
    link.classList.add("active-channel");
};


// A reusable function to apply the listeners to all channel links
function bindChannelListeners() {
    console.log("Binding channel link listeners...");
    document.querySelectorAll(".channel-link").forEach(link => {
        // Use a flag to prevent attaching the listener multiple times to the same element
        if (!link.classList.contains('listeners-bound')) {
            link.addEventListener("click", channelLinkClickHandler);
            link.classList.add('listeners-bound');
        }
    });
}

function wsConnect(action, address, password, username) {
// action is if ur signin in or registering
// msg is the object with the input field
// address is the object with the server address

    let url = "ws://" + address + "/ws"
    webSocket = new WebSocket(url);

    let sidebar = document.getElementById("sidebar")
    function inputInit() {
        console.log("input init")

        // Call the new function to bind listeners (initial bind)
        bindChannelListeners(); 

        // Find the main channel and visually mark it as active on first login
        const mainLink = document.querySelector('.channel-link[data-channel="main"]');
        if (mainLink) {
            mainLink.classList.add("active-channel");
            // NOTE: If you want 'main' to be connected on login, 
            // you would call connectChannel("main") here or inside loginStats
        }

        // Your existing message input listener
        msgInput.addEventListener("keydown", function(e) {
            if (e.key === "Enter") {
                e.preventDefault();
                writeMessage();
                messageInput.value = "";
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

        console.log("--- server: ", event.data)
        const msg = JSON.parse(event.data);

        switch (msg.Rtype) {
            /* --- Loggin In and Registering */
            case "loginStats":
                if (msg.Status == "ok") {
                    userToken = msg.Value

                    setDetails(msg.ServerName, msg.ServerDesc)
                    guiTransition()
                    inputInit()

                let payload = {
                    Rtype: "channelsList",
                    Username: username,
                    SessionToken: password,
                };
                console.log("Sent channels list req to server at", address)
                webSocket.send(JSON.stringify(payload));

                } else if (msg.Status == "invalidInput") {
                    showAlert("Invalid username: username must only contain normal letters (capital included), numbers and underscores. Dashes ('-') are not supported.")

                } else if (msg.Status == "alreadyLoggedIn") {
                    showAlert("Already logged in: A session already exists for this user")
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

            case "channelList": {
                console.log("new channel list: ", msg);

                let listEl = document.getElementById("channelList");

                // clear existing list EXCEPT for "main"
                listEl.innerHTML = `<li><a href="#" class="channel-link" data-channel="main">main</a></li>`;

                // loop channels from server and create new <li> elements
                for (const [name, count] of Object.entries(msg.ChannelList)) {
                    if (name === "main") continue;
                    let li = document.createElement("li");
                    li.innerHTML = `<a href="#" class="channel-link" data-channel="${name}">${name} (${count})</a>`;
                    listEl.appendChild(li);
                }
                
                // re-bind listeners to include the new channels
                bindChannelListeners(); 

                break;
            }

            case "kicked":
                showAlert("Connection force-closed by the server");
                logout();
                break;

            case "rejected":
                showAlert("Login was rejected by the server")
                break;

            /* --- Client stuff --- */
            case "disconnectStats":
                if (msg.Status == "error") {
                    showAlert("Error disconnecting " + msg.Value)
                } else {
                    console.log("disconnect ok")
                }
                break;

            case "invalidChannel":
                showAlert("Invalid Channel")
                console.warn("Invalid channel")
                break;

            case "unknownReq":
                console.warn("Outdated client or server, server sent unknownReq")
                break;

            case "messageCache":
                if (msg.MsgCache == null && msg.MsgCache == undefined) {

                } else {
                    for (const [username, message] of Object.entries(msg.MsgCache)) {
                        messageDisplay(username, message);
                    }
                }
                break;

            case "isr":
                showAlert("[Server] Internal Server Error: " + msg.Status)
                logout()
                break;
            

            case "invalidToken":
                showAlert("An invalid token was provided to the server -- please login again.")
                break;
                
            case "alreadyConnected":
                showAlert("You are already connected to a channel")
                break;

            case "connectStats":
                if (msg.Status == "error") {
                    showAlert("Error disconnecting " + msg.Value)
                } else if (msg.Status == "invalidToken") {
                    showAlert("Invalid token: please reload this page and login again, if the issue persists, contact an administrator.")
                } else {
                    console.log("connect OK")
                }
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
                messageDisplay(msg.Username, msg.Value, msg.Time)
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
    usernameV = username.value
    let payload = {
                Rtype: "disconnect",
                SessionToken: userToken,
                ChannelID: currentChannel,
                Username: usernameV,
    };
    console.log("Disconnecting: ", payload)
    webSocket.send(JSON.stringify(payload));
}

function connectChannel(channel) {
    alreadyRan = true
    msgValue = msgInput.value
    usernameV = username.value
    let payload = {
                Rtype: "connect",
                SessionToken: userToken,
                ChannelID: channel,
                Username: usernameV,
    };
    currentChannel = channel
    console.log("Connected to channel: ", payload)
    webSocket.send(JSON.stringify(payload));
}

function writeMessage() {
    aUsername = document.getElementById("username").value
    msgValue = msgInput.value
    let payload = {
                Rtype: "message",
                SessionToken: userToken,
                ChannelID: currentChannel,
                Message: msgValue,
                Username: aUsername,
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

function messageDisplay(username, message, time) {
    const messagesDiv = document.getElementById("messages");

    const wrapper = document.createElement("div");

    // Username
    const userEl = document.createElement("b");
    userEl.textContent = username + ": ";

    // Message
    const msgEl = document.createElement("span");
    msgEl.innerHTML = formatMessage(message); // safe + markdown

    // Timestamp
    const timeEl = document.createElement("span");
    timeEl.style.color = "gray"; // make it gray
    timeEl.style.marginLeft = "6px"; // optional spacing
    const date = new Date(time);

    if (isNaN(date.getTime())) {
        // Invalid date
        timeEl.textContent = "No Time";
        timeEl.title = "No Time";
    } else {
        // Format hh:mm for display
        const hours = String(date.getHours()).padStart(2, "0");
        const minutes = String(date.getMinutes()).padStart(2, "0");
        timeEl.textContent = `${hours}:${minutes}`;

        // Full human-readable tooltip: "July 3rd 2025 at 14:23:45"
        const day = date.getDate();
        const daySuffix = (d) => {
            if (d > 3 && d < 21) return "th";
            switch (d % 10) {
                case 1: return "st";
                case 2: return "nd";
                case 3: return "rd";
                default: return "th";
            }
        };
        const monthNames = [
            "January", "February", "March", "April", "May", "June",
            "July", "August", "September", "October", "November", "December"
        ];
        const fullTime = `${monthNames[date.getMonth()]} ${day}${daySuffix(day)} ${date.getFullYear()} at ${hours}:${minutes}:${String(date.getSeconds()).padStart(2,"0")}`;
        timeEl.title = fullTime; // tooltip
    }

    // Append elements
    wrapper.appendChild(userEl);
    wrapper.appendChild(msgEl);
    wrapper.appendChild(timeEl);

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

let userpopupBox = document.getElementById("user-popup")
let Pusername = document.getElementById("user-popup-username")
let Ptoken = document.getElementById("user-popup-token")
let Prank = document.getElementById("user-rank")

function showUserPopup() {
    if (username.value == "admin" ) {
        Pusername.innerHTML = username.value
        Pusername.style.color = "gold"
        Pusername.style.cursor = "crosshair"
    } else {
        Pusername.innerHTML = username.value
    }
    Ptoken.innerHTML = userToken
    Ptoken.title = "Click to copy!";
    //Prank.innerHTML = userRank
    userpopupBox.style.display = 'block'
}

// Make it clickable
Ptoken.style.cursor = "pointer"; // optional, shows pointer on hover
Ptoken.onclick = () => {
    const tokenText = Ptoken.textContent || Ptoken.innerText;
    if (!tokenText) return;

    // Copy to clipboard
    navigator.clipboard.writeText(tokenText)
        .then(() => {
            alert("Token copied to clipboard!");
        })
        .catch(err => {
            console.error("Failed to copy token:", err);
        });
};


let usernameLink = document.getElementById("usernamebox")
usernameLink.onclick = () => {
    showUserPopup()

    return false;
};

let currentStatus = ""; // variable to hold the value
let statusForm = document.getElementById("statusForm")

statusForm.onsubmit = (e) => {
    e.preventDefault(); // stop page reload

    const statusSelect = document.getElementById("statuses");
    currentStatus = statusSelect.value;

    console.log("Status set to:", currentStatus);
};
