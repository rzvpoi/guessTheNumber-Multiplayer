<script setup>
import { onMounted, ref, watch} from 'vue'
import { useRoute} from 'vue-router';
import { receivedMessage, sendAction} from '../webSocket';

const route = useRoute();
const code_param= route.params.code;
const code = code_param.split('=')[1]
const player1_param = route.params.leader
const player1 = player1_param.split('=')[1]
const player2_param = route.params.player
const player2 = player2_param.split('=')[1]

const messageReceived = ref(false);
let isPlayerNotJoin = player2.length ? false : ref(true)

let startGame = ref(false)

onMounted( ()=> {
        
        document.getElementById('player1_username').innerText = player1
        document.getElementById('room_code').innerText = code
        if (player2.length > 0) {
            console.log("Player 2 captured")
            document.getElementById('player2_username').innerText = player2
        } 
})


watch(receivedMessage, (newVal) => {
  messageReceived.value = !!newVal; // Update messageReceived based on receivedMessage value
  const msg = receivedMessage._rawValue
  
  if(msg.event == "player_joined") {
    isPlayerNotJoin.value = false
    startGame = true
    document.getElementById('player2_username').innerText = msg.player
  }  
});

function startTheGame() { 
    sendAction({action: 'start_game', code: code, player1: player1, player2: player2})
}

</script>

<template>
     <div class="container lobby-container">
        <h1 class="text-center mb-4">Lobby Waiting Room</h1>
        <div class="row">
          <div class="col-md-6">
            <h2>Players in lobby</h2>
            <ul class="waiting-list">
              <li class="waiting-list-item">
                <p id="player1_username"></p>
              </li>
              <li class="waiting-list-item">
                <p id="player2_username"></p>   
                <div v-if="isPlayerNotJoin == true">
                    <div class="d-flex align-items-center">
                        <strong>Waiting for another player...   </strong>
                        <div class="spinner-border ms-auto spinner-border-sm" role="status" aria-hidden="true"></div>
                    </div>
                </div>
              </li>
            </ul>
            <div class="mb-1">
                <div v-if="startGame">
                    <button class="btn btn-primary" @click="startTheGame">Start Game</button>
                </div>
            </div>
          </div>
          <div class="col-md-6">
            <h2>Room Code</h2>
            <div class="card">
              <div class="card-body">
                <h5 class="card-title" id="room_code">N4H1RR</h5>
                <p class="card-text">Share this code with your friend so they can join!</p>
              </div>
            </div>
            <!-- Add more announcements as needed -->
          </div>
        </div>
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
  