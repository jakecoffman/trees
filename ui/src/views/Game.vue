<template>
  <section>
    <game-header :game="game" :you="you"/>
    <hex-grid v-if="game" :game="game" :you="you" :selection="selection" class="hex-grid" @select="select"/>
    <game-footer ref="foot" v-if="game" :game="game" :selection="selection" :you="you"></game-footer>

    <div class="smokescreen" v-if="!game || game.Players.length !== 2 || conn !== 'Open'"></div>
    <modal :open="game && game.Players.length !== 2">
      <p>Waiting for opponent</p>
      <p>Room code {{game.Code}}</p>
    </modal>
    <modal :open="conn !== 'Open'">
      <p>{{conn}}</p>
      <button @click="$router.push('/')">Go home</button>
    </modal>
  </section>
</template>
<script>
import HexGrid from "../components/HexGrid.vue";
import Modal from "../components/Modal.vue";
import GameHeader from "../components/GameHeader.vue";
import { createToast } from 'mosha-vue-toastify';
import 'mosha-vue-toastify/dist/style.css'
import GameFooter from "../components/GameFooter.vue";
import {computed} from 'vue'

export default {
  components: {
    GameFooter,
    GameHeader,
    Modal,
    HexGrid,
  },
  data() {
    return {
      ws: null,
      conn: 'Connecting...',
      game: null,
      you: null,
      showFooter: false,
      selection: null
    }
  },
  provide() {
    return {
      ws: computed(() => this.ws),
      conn: computed(() => this.conn)
    }
  },
  async mounted() {
    try {
      const r = await fetch('/api/login', {credentials: "include"})
      if (!r.ok) {
        return alert(await r.text())
      }
      await r.json()
    } catch (e) {
      console.error(e)
      return
    }

    let protocol = 'ws://'
    if (location.protocol === 'https:') {
      protocol = 'wss://'
    }
    this.ws = new WebSocket(`${protocol}${location.host}/ws?code=${this.$route.params.id}`)
    this.ws.onopen = this.wsOpen
    this.ws.onclose = this.wsClose
    this.ws.onerror = this.wsError
    this.ws.onmessage = this.wsMessage
  },
  beforeUnmount() {
    try {
      this.ws.close()
    } finally {
      this.ws = null
    }
  },
  methods: {
    wsOpen() {
      this.conn = 'Open'
    },
    wsClose() {
      this.conn = 'Connection closed'
    },
    wsError() {
      this.conn = 'Error on connection'
      createToast('WebSocket error', {type: 'danger', position: 'bottom-right'})
    },
    wsMessage(msg) {
      const data = JSON.parse(msg.data)
      switch (data.Kind) {
        case "msg":
          createToast(data.Value, {type: 'danger', position: 'bottom-right'})
          break
        case "unlock":
          this.$refs.foot.unlock()
          break
        case "room":
          console.log(data.Room)
          if (this.$route.params.id !== data.Room.Code) {
            this.$router.replace(`/games/${data.Room.Code}`)
          }
          this.game = data.Room
          this.you = data.You
          console.log("You", this.you)
          break
        default:
          alert('unhandled message:' + msg.data)
          break
      }
    },
    select(index) {
      this.showFooter = true
      this.selection = index
    }
  }
}
</script>
<style scoped>
.smokescreen {
  position: fixed;
  z-index: 1;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  opacity: 50%;
  background: black;
}
</style>
