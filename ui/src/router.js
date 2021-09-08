import {createRouter, createWebHistory} from 'vue-router'

import Home from './views/Home.vue'
import NewGame from './views/NewGame.vue'
import Game from './views/Game.vue'
import Rules from './views/Rules.vue'
import NotFound from './views/NotFound.vue'
import Tutorial from './views/Tutorial.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {path: '/', component: Home, name: 'Home'},
    {path: '/games', component: NewGame, name: 'New Game'},
    {path: '/games/:id', component: Game, name: 'Game'},
    {path: '/rules', component: Rules, name: 'Rules'},
    {path: '/tutorial', component: Tutorial, name: 'Tutorial'},
    {path: '/:catchAll(.*)', component: NotFound, name: 'Not Found'}
  ]
})
