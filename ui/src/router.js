import {createRouter, createWebHistory} from 'vue-router'

import Home from './views/Home.vue'
import Game from './views/Game.vue'
import Rules from './views/Rules.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {path: '/', component: Home, name: 'Home'},
    {path: '/game', component: Game, name: 'New Game'},
    {path: '/game/:id', component: Game, name: 'Game'},
    {path: '/rules', component: Rules, name: 'Rules'},
  ]
})
