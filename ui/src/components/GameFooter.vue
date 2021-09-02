<template>
  <footer>
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
        <button v-if="tree.Size >= 1" @click="seed(selection)">Seed</button>
        <button v-if="tree.Size < 3" @click="grow(selection)">Grow</button>
        <button v-if="tree.Size === 3" @click="sell(selection)">Sell</button>
      </span>
    </div>
  </footer>
</template>
<script>
export default {
  props: ['game', 'selection', 'you'],
  inject: ['ws'],
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
    seed(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'seed', Index: selection}))
    },
    sell(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'sell', Index: selection}))
    },
    grow(selection) {
      this.ws.value.send(JSON.stringify({Kind: 'grow', Index: selection}))
    }
  }
}
</script>
<style scoped>
footer {
  position: fixed;
  bottom: 0;
  width: 100vw;
  height: 4rem;
  background: black;

  display: flex;
  gap: .5rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
