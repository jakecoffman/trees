<template>
  <section>
    <game-header :game="game"/>
    <hex-grid v-if="game" :game="game" class="hex-grid"/>
    <footer>
      <button>End Turn</button>
    </footer>

    <modal v-if="game && game.Players.length !== 2">
      <p>Waiting for opponent</p>
      <p>Send them this link: <a :href="url">share</a></p>
    </modal>
  </section>
</template>
<script>
import HexGrid from "../components/HexGrid.vue";
import Modal from "../components/Modal.vue";
import GameHeader from "../components/GameHeader.vue";
import { createToast } from 'mosha-vue-toastify';
import 'mosha-vue-toastify/dist/style.css'

export default {
  components: {
    GameHeader,
    Modal,
    HexGrid,
  },
  computed: {
    url() {
      return location.href
    }
  },
  data() {
    return {
      ws: null,
      state: 'Connecting',
      game: null
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

    let params = ''
    if (this.$route.params.id) {
      console.log('JOINING')
      params = `action=join&code=${this.$route.params.id}`
    } else {
      console.log('NEWING')
      params = `action=new`
    }

    if (!this.ws) {
      this.ws = new WebSocket(`ws://${location.host}/ws?${params}`)
      this.ws.onopen = this.wsOpen
      this.ws.onclose = this.wsClose
      this.ws.onerror = this.wsError
      this.ws.onmessage = this.wsMessage
    }
  },
  methods: {
    wsOpen() {
      this.state = 'Open'
      console.log('open')
    },
    wsClose() {
      this.state = 'Closed'
      createToast('Disconnected, please refresh', {type: 'danger', position: 'bottom-right'})
    },
    wsError() {
      this.state = 'Error'
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
          if (!this.$route.params.id) {
            this.$router.replace(`/game/${data.Room.Code}`)
          }
          this.game = data.Room
          break
        default:
          alert('unhandled message:' + msg.data)
          break
      }
    }
  }
}
</script>
<style scoped>
.hex-grid {
  padding-bottom: 3rem;
}
footer {
  position: fixed;
  bottom: 0;
  width: 100vw;
  height: 3rem;
  background: black;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
