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
