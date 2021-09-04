<template>
  <section>
    <game-header :game="game" :you="you"/>
    <hex-grid v-if="game" :game="game" :you="you" :selection="selection" class="hex-grid" @select="select"/>
    <game-footer v-if="game" :game="game" :selection="selection" :you="you"></game-footer>
    <modal v-if="game && game.Players.length !== 2">
      <p>Waiting for opponent</p>
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
      conn: 'Connecting',
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
      console.log('open')
    },
    wsClose() {
      this.conn = 'Closed'
    },
    wsError() {
      this.conn = 'Error'
      createToast('WebSocket error', {type: 'danger', position: 'bottom-right'})
    },
    wsMessage(msg) {
      const data = JSON.parse(msg.data)
      switch (data.Kind) {
        case "msg":
          createToast(data.Value, {type: 'danger', position: 'bottom-right'})
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
</style>
