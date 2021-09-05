<template>
  <header v-if="game">
    <div class="score-mini">
      <div class="orange">
        {{ you === 0 ? 'You' : 'Opponent' }}: {{game.State.Score[0]}} ☀️{{ game.State.Energy[0] }}
      </div>
      <div class="blue">
        {{ you === 1 ? 'You' : 'Opponent' }}: {{game.State.Score[1]}} ☀️{{ game.State.Energy[1] }}
      </div>
    </div>
    <div class="score">
      <fieldset class="orange">
        <legend>{{ you === 0 ? 'You' : 'Opponent' }}</legend>
        Score {{ game.State.Score[0] }}
        <span v-if="!game.Players[0].Connected">(Disconnected)</span>
        <br/>
        ☀️{{ game.State.Energy[0] }}
      </fieldset>

      <fieldset class="blue">
        <legend>{{ you === 1 ? 'You' : 'Opponent' }}</legend>
        <div v-if="game.Players.length < 2">
          Waiting
        </div>
        <div v-else>
          Score {{ game.State.Score[1] }}
          <span v-if="!game.Players[1].Connected">(Disconnected)</span>
          <br/>
          ☀️{{ game.State.Energy[1] }}
        </div>
      </fieldset>
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
fieldset {
  border: 1px solid;
  border-radius: 5px;
}
.score-mini {
  display: none;
}
@media (max-height: 600px) {
  .score {
    display: none;
  }
  .score-mini {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
