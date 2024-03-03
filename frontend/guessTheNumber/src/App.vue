<script setup>
import { ref, computed, watch, onBeforeUnmount  } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import Lobby from './views/Lobby.vue'
import Menu from './views/Menu.vue'
import Game from './views/Game.vue'
import { useWebSocket} from './webSocket';
import { receivedMessage } from './webSocket';
import router from './router/index.js'

const { connected } = useWebSocket();

const routes = {
  '/': Menu,
  '/lobby': Lobby,
  "/game": Game
}

const messageReceived = ref(false);

watch(receivedMessage, (newVal) => {
  messageReceived.value = !!newVal; // Update messageReceived based on receivedMessage value
  const msg = receivedMessage._rawValue
  
  if(msg.event == "lobby_created") {
    const attributes = `code=${msg.code}/leader=${msg.player1}/player=`
    router.push({path: `/lobby/${attributes}`})

  } else if(msg.event == "join_lobby") {
    const attributes = `code=${msg.code}/leader=${msg.leader}/player=${msg.player}`
    router.push({path: `/lobby/${attributes}`})
  } else if( msg.event == "start_game") {
    const attributes = `code=${msg.code}/player=${msg.player}/number=${msg.player_number}`
    router.push({path: `/game/${attributes}`})
  }
});
const confirmBeforeRefresh = (event) => {
  // Display confirmation message
  const confirmationMessage = 'Are you sure you want to leave? Your connection to the WebSocket will be lost.';

  // Display confirmation dialog
  event.preventDefault();
  event.returnValue = confirmationMessage;
};

// Add event listener when component is mounted
window.addEventListener('beforeunload', confirmBeforeRefresh);
</script>

<template>
  <div v-if="connected">
    <RouterView />
  </div>
  <div v-else>
    <p>WebSocket is not connected.</p>
  </div>

 
</template>

<style scoped>
</style>