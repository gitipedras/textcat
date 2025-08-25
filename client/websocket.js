let ws;

    function connect() {
        ws = new WebSocket("ws://localhost:8080/ws");

        ws.onopen = function() {
            console.log("Connected to WebSocket server");
        };

        ws.onmessage = function(event) {
            let messageDisplay = document.getElementById("messages");
            messageDisplay.innerHTML += `<p>${event.data}</p>`;
        };

        ws.onclose = function() {
            console.log("WebSocket connection closed, retrying...");
            setTimeout(connect, 200 ); // Reconnect after 1 second
        };

        ws.onerror = function(error) {
            console.error("WebSocket error:", error);
        };
    }

    function sendMessage() {
        let input = document.getElementById("messageInput");
        let message = input.value;
        ws.send(message);
        input.value = "";
    }

    connect();