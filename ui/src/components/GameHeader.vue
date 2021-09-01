<template>
  <header class="score" v-if="game">
      <span class="orange">
        <span>Player 1 Score {{game.State.Score[0]}} Energy {{game.State.Energy[0]}}</span>
        <span v-if="you === 0">(you)</span>
        <span v-if="!game.Players[0].Connected">(Disconnected)</span>
      </span>
    <span class="blue">
      <span v-if="game.Players.length < 2">Waiting for opponent</span>
      <span v-else>
        <span>Player 2 Score {{game.State.Score[1]}} Energy {{game.State.Energy[1]}}</span>
        <span v-if="you === 1">(you)</span>
        <span v-if="!game.Players[1].Connected">(Disconnected)</span>
      </span>
      </span>
    <span>OPT</span>
    <span class="flex">
      <span class="grow">Day {{game.State.Day}} Nutrients {{game.State.Nutrients}}</span>
      <span v-if="conn !== 'Open'">Connection {{conn}}</span>
    </span>
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
  grid-template-columns: 1fr 1fr 3rem;
  background: black;
  color: white;
}
</style>
