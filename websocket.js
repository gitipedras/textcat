var userToken
let currentChannel = "main";

function wsConnect(action, msgInput, address, password, username) {
// action is if ur signin in or registering
// msg is the object with the input field
// address is the object with the server address

    let url = "ws://" + address + "/ws"
    webSocket = new WebSocket(url);

    function inputInit() {
      const bind = () => {
        const input = document.getElementById(msgInput);
        if (!input) {
          console.warn("Input not found:", msgInput);
          return;
        }

        input.addEventListener("keydown", function (event) {
          if (event.key === "Enter") {
            if (!userToken) {
              showAlert("You must be logged in before sending messages.");
              return;
            }
            writeMessage(token, this.value);
            this.value = "";
          }
        });

        document.querySelectorAll(".channel-link").forEach(link => {
          link.addEventListener("click", e => {
            e.preventDefault();
            currentChannel = link.dataset.channel;
            console.log("Switched to channel:", currentChannel);
          });
        });
      };

      if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", bind, { once: true });
      } else {
        bind();
      }
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

    console.log(password)

    webSocket.onmessage = function(event) {

        const msg = JSON.parse(event.data);
        console.log("message is ", msg)

        switch (msg.Rtype) {
            /* --- Loggin In and Registering */
            case "loginStats":
                if (msg.Status == "ok") {
                    console.log(msg.Value)
                    userToken = msg.Value
                    guiTransition()
                    inputInit()
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

            case "kicked":
                showAlert("Connection force-closed by the server");
                logout();
                break;

            case "rejected":
                showAlert("Login was rejected by the server")
                break;

            /* --- Client stuff --- */

            case "invalidChannel":
                showAlert("Invalid Channel: " + msg.Value)
                console.warn("Invalid channel: " + msg.Value)
                break;

            case "unknownReq":
                console.warn("Outdated client or server, server sent unknownReq")
                break;

            /*case "isr":
                showAlert("[Server] Internal Server Error")
                logout()
                break;
            */

            case "alreadyConnected":
                showAlert("You are already connected to this channel!")
                break;

            /*
            case "invalidSession":
                showAlert("An invalid session was provided")
                logout()
                break;
            */

            case "NewMessage":
                console.log("New message: " + msg.Value)
                messageDisplay.innerHTML += `<p>@${msg.username}: ${msg.message}</p>`;
                break;

            default:
                console.log("Unknown request type from server: " + msg.Rtype);
                break;
        }

    };

    webSocket.onclose = function() {
        console.log("WebSocket connection closed.");
        //showAlert("Connection closed. Reason: unknown")
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

function writeMessage(token, messagecontent, msgUser) {
    let payload = {
                Rtype: "message",
                SessionToken: userToken,
                ChannelID: currentChannel,
                Message: messagecontent,
                Username: msgUser,
    };
    console.log("Sent message: ", payload)
    webSocket.send(JSON.stringify(payload));
}