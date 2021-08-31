<template>
  <section>
    <hex-grid/>
    State: {{ state }}
  </section>
</template>
<script>
import HexGrid from "../components/HexGrid.vue";

export default {
  components: {HexGrid},
  data() {
    return {
      ws: null,
      state: 'Connecting'
    }
  },
  async mounted() {
    await fetch('http://127.0.0.1:8080/login', {credentials: 'include'})

    this.ws = new WebSocket('ws://127.0.0.1:8080/ws')
    this.ws.onopen = this.wsOpen
    this.ws.onclose = this.wsClose
    this.ws.onerror = this.wsError
    this.ws.onmessage = this.wsMessage
  },
  methods: {
    wsOpen() {
      this.state = 'Open'
    },
    wsClose() {
      this.state = 'Closed'
    },
    wsError() {
      this.state = 'Error'
    },
    wsMessage(msg) {
      console.log(msg)
      const data = JSON.parse(msg.data)
      switch (data.Kind) {
        default:
          alert('unhandled message:' + msg.data)
          break
      }
    }
  }
}
</script>
