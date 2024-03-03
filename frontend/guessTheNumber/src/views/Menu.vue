<script setup>
import { sendAction as sendAction } from '../webSocket';
import { useRoute } from 'vue-router';

const route = useRoute();
const createNewGame = () => {
  let username = document.getElementById('usernameInput').value
    if(username==""){
        document.getElementById('errorName').innerText = "Username cannot be empty!"
        return
    }
    sendAction({ action: 'create_lobby', username:username});
}


const joinGame = () => {
  let code = document.getElementById('codeInput').value
  let username = document.getElementById('usernameJoinInput').value

  if(code==""){
        document.getElementById('errorCode-Join').innerText = "Code cannot be empty!"
        return
    } else {
      document.getElementById('errorCode-Join').style.display = "none"
    }

  if(username==""){
        document.getElementById('errorName-Join').innerText = "Username cannot be empty!"
        return
    } else {
      document.getElementById('errorCode-Join').style.display = "none"
    }
  
  sendAction({ action: 'join_lobby', code:code, username:username});
}

</script>

<template>
  

  <div class="row justify-content-center">
  <div>
    <h1 class="text-center mb-4">Guess the number! Multiplayer</h1>
  </div>
  <div class="col-md-4 mb-3">
    <button type="button" class="btn btn-primary btn-lg btn-block" data-bs-toggle="modal" data-bs-target="#createLobby" >
      New Game
    </button>
  </div>
  <div class="col-md-4">
    <button type="button" class="btn btn-primary btn-lg btn-block" data-bs-toggle="modal" data-bs-target="#joinLobby">
      Join Game
    </button>
  </div>
</div>


<!-- Create a new game Modal -->
<div class="modal fade" id="createLobby" tabindex="-1" aria-labelledby="createLobbyLabel" aria-hidden="true" >
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="createLobbyLabel">Create a new game</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
        <form>
          <div class="form-group">
            <label for="usernameInput">Username</label>
            <input type="text" class="form-control" id="usernameInput" aria-describedby="usernameHelp" placeholder="Be creative!">
            <small id="errorName" class="form-text"></small>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary" @click="createNewGame" data-bs-dismiss="modal">Create</button>
      </div>
    </div>
  </div>
</div>

<!-- Join a game Modal -->
<div class="modal fade" id="joinLobby" tabindex="-1" aria-labelledby="joinLobbyLabel" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="joinLobbyLabel">Join a game</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
        <form>
          <div class="form-group">
            <label for="codeInput">Code</label>
            <input type="text" class="form-control" id="codeInput" placeholder="Get the lobby code from your friend">
            <small id="errorCode-Join" class="form-text"></small>
          </div>
          <div class="form-group">
            <label for="usernameInput">Username</label>
            <input type="text" class="form-control" id="usernameJoinInput" aria-describedby="usernameHelp" placeholder="Be creative!">
            <small id="errorName-Join" class="form-text"></small>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary" @click="joinGame" data-bs-dismiss="modal">Join</button>
      </div>
    </div>
  </div>
</div>

<div>
  <p> ATTENTION: This game is still in beta. Do not try to refresh the page because the connection will be lost!</p>
</div>

</template>
  
  <style>
  @media (min-width: 1024px) {
    .about {
      min-height: 100vh;
      display: flex;
      align-items: center;
    }
  }
  </style>
  