<template>
  <article>
    Creating new game...
  </article>
</template>
<script>
export default {
  async mounted() {
    const r = await fetch('/api/rooms', {
      credentials: 'include',
      method: 'POST'
    })
    if (!r.ok) {
      return alert("Failed to create game: " + await r.text())
    }
    const data = await r.json()
    await this.$router.replace(`/games/${data.Code}`)
  }
}
</script>
