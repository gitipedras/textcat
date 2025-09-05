const base64Cursor = "data:image/png;base64,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";
const serverInput = document.getElementById("server")
serverInput.value = defaultServer;

// Apply custom cursor to the whole page
document.body.style.cursor = `url(${base64Cursor}) 0 0, auto`;

/* gui */
let logingui = document.getElementById("logingui")
let settingsui = document.getElementById("settingsPage")
let msgui = document.getElementById("maingui")
let userbox = document.getElementById("userBox")
let chatbar = document.getElementById("chatInputBar")

/* =================== GUI ======================= */
/* ----------------------------------------------- */

/* default gui states */
logingui.style.display = "block"
maingui.style.display = "none"
sidebar.style.display = "none"
settingsui.style.display = "none"
msgui.style.display = "none"
chatbar.style.display = "none"
userbox.style.display = "none"

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

/* =================== APP FUNCTIONS ======================= */
/* --------------------------------------------------------- */

let loggingIn
let cStats = document.getElementById("cStats")

function login() {
    cStats.innerHTML = 'Connecting to server...'
	loggingIn = true
	startWebsocket(loggingIn)
}

function loginRegister() {
    cStats.innerHTML = 'Connecting to server...'
	loggingIn = false
	startWebsocket(loggingIn)
}

function deleteAllStorage() {
    localStorage.clear();
    alert("Deleted localStorage.")
}

function showAlert(message) {
  const popup = document.getElementById('custom-popup');
  const popupText = document.getElementById('popup-text');

  popupText.textContent = message;
  popup.style.display = 'flex';
}

function closePopup() {
  const popup = document.getElementById('custom-popup');
  popup.style.display = 'none';
}


function showUser(user) {
  const popup = document.getElementById('user-popup');
  document.getElementById('user-popup-username').textContent = user.username;
  document.getElementById('user-popup-description').textContent = "Description: " + (user.description || "No description");
  document.getElementById('user-popup-date').textContent = "Date Created: " + (user.dateCreated || "Unknown");
  document.getElementById('user-popup-token').textContent = "" + (user.token || "N/A");

  popup.style.display = 'flex';
}

document.getElementById('user-popup-ok').onclick = function() {
  document.getElementById('user-popup').style.display = 'none';
};

/* =================== WEBSOCKETS ;) ======================= */
/* --------------------------------------------------------- */

function startWebsocket(loggingIn) {
    let wsServer = document.getElementById("server").value
    let password = document.getElementById("password").value;
    let username = document.getElementById("username").value
	if (loggingIn == true) {
		wsConnect("login", wsServer, password, username)
	} else {
		wsConnect("register", wsServer, password, username)
	}
}

document.addEventListener("DOMContentLoaded", () => {
    const themeSelect = document.getElementById("themes");
    const form = document.getElementById("themeForm");

    // Load saved theme
    const savedTheme = localStorage.getItem("theme") || "light";
    themeSelect.value = savedTheme;
    applyTheme(savedTheme);

    // Handle form submit
    form.addEventListener("submit", (e) => {
        e.preventDefault();
        const selectedTheme = themeSelect.value;
        localStorage.setItem("theme", selectedTheme);
        applyTheme(selectedTheme);
    });

    // Function to apply theme
    function applyTheme(theme) {
        document.body.setAttribute("data-theme", theme);
    }
});