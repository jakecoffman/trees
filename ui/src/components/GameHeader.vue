<template>
  <header v-if="game">
    <div class="score">
      <span class="orange">
        <span>
          {{ you === 0 ? 'You' : 'Player 1' }}:
          <span>Score {{ game.State.Score[0] }}<br/>☀️{{ game.State.Energy[0] }}</span>
          <span v-if="!game.Players[0].Connected">(Disconnected)</span>
        </span>
      </span>

      <span class="blue">
        <span v-if="game.Players.length < 2">Need opponent</span>
        <span v-else>
          <span>
            {{ you === 1 ? 'You' : 'Player 2' }}:
            <span>Score {{ game.State.Score[1] }}<br/>☀️{{ game.State.Energy[1] }}</span>
            <span v-if="!game.Players[1].Connected">(Disconnected)</span>
          </span>
        </span>
      </span>
    </div>

    <div class="infos">
      <span>Day {{game.State.Day}}</span>
      <span>Nutrients {{game.State.Nutrients}}</span>
      <span>Room {{game.Code}}</span>
    </div>

  </header>
</template>
<script>
export default {
  props: ['game', 'conn', 'you']
}
</script>
<style scoped>
.score {
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: black;
  color: white;
}
.infos {
  margin-top: .5rem;
  margin-bottom: .5rem;
  display: flex;
  justify-content: space-between;
}
</style>
