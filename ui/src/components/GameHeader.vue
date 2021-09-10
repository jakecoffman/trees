<template>
  <header v-if="game">
    <div class="score-mini">
      <score-board :you="you" :score="game.State.Score[you]" :energy="game.State.Energy[you]"/>
      <score-board :you="you" :score="game.State.Score[you]" :energy="game.State.Energy[you]"/>
    </div>
    <div class="score">
      <fieldset class="orange">
        <legend>{{ you === 0 ? 'You' : 'Opponent' }}</legend>
        <high-light>
          <span :key="game.State.Score[0]">Score {{ game.State.Score[0] }}</span>
        </high-light>
        <high-light>
          <span v-if="!game.Players[0].Connected">(Disconnected)</span>
        </high-light>
        <br/>
        <high-light>
          <span :key="game.State.Energy[0]">☀️{{ game.State.Energy[0] }}</span>
        </high-light>
      </fieldset>

      <fieldset class="blue">
        <legend>{{ you === 1 ? 'You' : 'Opponent' }}</legend>
        <div v-if="game.Players.length < 2">
          Waiting
        </div>
        <div v-else>
          <high-light>
            <span :key="game.State.Score[1]">Score {{ game.State.Score[1] }}</span>
          </high-light>
          <high-light>
            <span v-if="!game.Players[1].Connected">(Disconnected)</span>
          </high-light>
          <br/>
          <high-light>
            <span :key="game.State.Energy[1]">☀️{{ game.State.Energy[1] }}</span>
          </high-light>
        </div>
      </fieldset>
    </div>

    <div class="infos">
      <high-light>
        <span :key="game.State.Day">Day {{game.State.Day}}</span>
      </high-light>
      <high-light>
        <span :key="game.State.Nutrients">Nutrients {{game.State.Nutrients}}</span>
      </high-light>
      <span>Room {{game.Code}}</span>
    </div>

  </header>
</template>
<script>
import ScoreBoard from "./ScoreBoard.vue";
import HighLight from "./HighLight.vue";
export default {
  components: {HighLight, ScoreBoard},
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
