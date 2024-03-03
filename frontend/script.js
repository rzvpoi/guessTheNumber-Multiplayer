
// Function to handle incoming messages from the WebSocket server
function handleMessage(event) {
    // Handle incoming messages from the WebSocket server
    const message = JSON.parse(event.data);
    // Update UI based on the received message
    console.log(message)

    if(message.response == "create_lobby") {
        console.log("Lobby code is "+ message.code)
        location.href = location.origin+"/frontend/lobby.html?code="+message.code+"&player1="+message.player1;
    }
    
}

// Function to send player actions to the WebSocket server
function sendAction(action) {
    // Send player actions to the WebSocket server
    socket.send(JSON.stringify(action));
}

// Create a WebSocket connection
const socket = new WebSocket('ws://localhost:8080/ws');

// Event listener for WebSocket connection opened
socket.addEventListener('open', () => {
    console.log('WebSocket connection opened');
});

// Event listener for WebSocket messages
socket.addEventListener('message', handleMessage);

// Event listener for WebSocket connection closed
socket.addEventListener('close', () => {
    console.log('WebSocket connection closed');
});

// Event listener for WebSocket connection errors
socket.addEventListener('error', (error) => {
    console.error('WebSocket error:', error);
});

//Event listener for button actions
document.getElementById('create_lobby').addEventListener('click', () => {
    let username = document.getElementById('usernameInput').value
    if(username==""){
        console.log("username empty")
        document.getElementById('error').innerText = "Username cannot be empty!"
        return
    }
    sendAction({ action: 'create_lobby', username:username});
});

document.getElementById('join_lobby').addEventListener('click', () => {
    let code = document.getElementById('codeInput').value
    let username = document.getElementById('usernameJoinInput').value
    sendAction({ action: 'join_lobby', code:code, username:username});
});

