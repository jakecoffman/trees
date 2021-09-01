<template>
  <section>
    <game-header :game="game"/>
    <hex-grid v-if="game" :game="game"/>
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
      params = `action=join&code=${this.$route.params.id}`
    } else {
      params = `action=new`
    }

    this.ws = new WebSocket(`ws://${location.host}/ws?${params}`)
    this.ws.onopen = this.wsOpen
    this.ws.onclose = this.wsClose
    this.ws.onerror = this.wsError
    this.ws.onmessage = this.wsMessage
  },
  methods: {
    wsOpen() {
      this.state = 'Open'
      console.log('open')
    },
    wsClose() {
      this.state = 'Closed'
    },
    wsError() {
      this.state = 'Error'
    },
    wsMessage(msg) {
      const data = JSON.parse(msg.data)
      console.log("got message", data.Kind)
      switch (data.Kind) {
        case "game":
          console.log(data.Game)
          if (!this.$route.params.id) {
            this.$router.replace(`/game/${data.Game.Code}`)
          }
          this.game = data.Game
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
