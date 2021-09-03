<template>
  <footer>
    Day {{game.State.Day}} Nutrients {{game.State.Nutrients}} Room {{game.Code}}
    <span v-if="conn.value !== 'Open'">Connection {{conn}}</span>

    <div v-if="selection">
      <span v-if="cell">Richness: {{cell.Richness}} &nbsp;</span>
      <span v-if="tree">
        <span>{{tree.Owner === you ? 'Your' : `Opponent's`}} Tree: Size {{tree.Size}}</span>
        <span v-if="tree.IsDormant">Dormant</span>
      </span>
    </div>
    <div>
      <button @click="endTurn()">End Turn</button>
      <span v-if="tree && tree.Owner === you && !tree.IsDormant">
        <button v-if="tree.Size >= 1" @click="seed1(selection)">Seed</button>
        <button v-if="tree.Size < 3" @click="grow(selection)">Grow</button>
        <button v-if="tree.Size === 3" @click="sell(selection)">Sell</button>
      </span>
      <span v-if="!tree && seedSource">
        <button @click="seed2(selection)">Seed Here</button>
      </span>
      <span v-if="seedSource">
        <button @click="seedSource = null">Cancel Seed</button>
      </span>
    </div>
  </footer>
</template>
<script>
export default {
  props: ['game', 'selection', 'you'],
  inject: ['ws', 'conn'],
  data() {
    return {
      seedSource: null
    }
  },
  computed: {
    cell() {
      if (!this.selection) {
        return null
      }
      return this.game.State.Board.Cells[this.selection]
    },
    tree() {
      if (!this.selection) {
        return null
      }
      return this.game.State.Trees[this.selection.toString()]
    }
  },
  methods: {
    endTurn() {
      this.ws.value.send(JSON.stringify({Kind: 'end'}))
    },
    seed1(selection) {
      this.seedSource = selection
    },
    seed2(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'seed', Source: this.seedSource, Target: selection}))
      this.seedSource = null
    },
    sell(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'sell', Target: selection}))
    },
    grow(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'grow', Target: selection}))
    }
  }
}
</script>
<style scoped>
footer {
  /*position: fixed;*/
  /*bottom: 0;*/
  /*width: 100vw;*/
  /*height: 4rem;*/
  /*background: black;*/

  display: flex;
  gap: .5rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
