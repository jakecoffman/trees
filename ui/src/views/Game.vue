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

    this.ws = new WebSocket(`ws://${location.host}/ws`)
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
