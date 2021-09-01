<template>
  <footer>
    Selected: {{selection}}
    Richness: {{game.State.Board.Cells[selection].Richness}}
    <span v-if="tree">
      <span>Tree: Owner {{tree.Owner === you ? 'you' : 'opponent'}} Size {{tree.Size}}</span>
      <span v-if="tree.IsDormant">Dormant</span>
      </span>
  </footer>
</template>
<script>
export default {
  props: ['game', 'selection', 'you'],
  computed: {
    tree() {
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
  height: 3rem;
  background: black;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
