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

function login() {
	loggingIn = true
	startWebsocket(loggingIn)
}

function loginRegister() {
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
    let password = document.getElementById("password").value

	if (loggingIn == true) {
		wsConnect("login", wsServer, password)
	} else {
		wsConnect("register", wsServer, password)
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