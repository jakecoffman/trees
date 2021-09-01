<template>
  <footer>
    <div v-if="selection">
      <span v-if="cell">Richness: {{cell.Richness}} &nbsp;</span>
      <span v-if="tree">
        <span>{{tree.Owner === you ? 'Your' : 'Opponent'}} Tree: Size {{tree.Size}}</span>
        <span v-if="tree.IsDormant">Dormant</span>
      </span>
    </div>
    <div>
      <button>End Turn</button>
      <span v-if="tree && tree.Owner === you && !tree.IsDormant">
        <button v-if="tree.Size >= 1">Seed</button>
        <button v-if="tree.Size < 3">Grow</button>
        <button v-if="tree.Size === 3">Sell</button>
      </span>
    </div>
  </footer>
</template>
<script>
export default {
  props: ['game', 'selection', 'you'],
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
