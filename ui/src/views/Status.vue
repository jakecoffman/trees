<template>
  <article>
    <h1>Status</h1>
    <div v-if="status">
      <p>Players Registered: {{status.PlayerCount}}</p>
      <p>Games Running: {{status.GameCount}}</p>
      <p>Active WS: {{status.ActiveWsConnections}}</p>
    </div>
  </article>
</template>
<script>
export default {
  data() {
    return {
      status: null
    }
  },
  mounted() {
    fetch(`/api/count`).then(r => {
      if (!r.ok) {
        return r.text().then(data => alert(data))
      }
      r.json().then(data => this.status = data)
    })
  }
}
</script>
<style scoped>
article {
  margin: 1rem;
}
</style>
