let socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = function(event) {
    console.log("Connected to server");
    sendMessage("join");  // Пример команды для присоединения к игре
};

socket.onmessage = function(event) {
    let messages = document.getElementById("messages");
    let message = document.createElement("div");
    message.textContent = event.data;
    messages.appendChild(message);
};

socket.onclose = function(event) {
    console.log("Disconnected from server");
};

socket.onerror = function(error) {
    console.error("WebSocket error:", error);
};

function sendMessage() {
    let input = document.getElementById("input");
    let message = input.value;
    socket.send(message);
    input.value = "";
}
