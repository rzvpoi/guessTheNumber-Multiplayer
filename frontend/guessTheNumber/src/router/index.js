import { createRouter, createWebHistory } from 'vue-router'
import Menu from '../views/Menu.vue'
import Lobby from '../views/Lobby.vue'
import Game from '../views/Game.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'menu',
      component: Menu
    },
    {
      path: '/lobby/:code/:leader/:player',
      name: 'lobby',
      component: Lobby
    },
    {
    path: '/game/:code/:player/:number',
    name: 'game',
    component: Game
    }
  ]
})

export default router
