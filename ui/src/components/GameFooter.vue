<template>
  <footer>
    Day {{game.State.Day}} Nutrients {{game.State.Nutrients}} Room {{game.Code}}
    <span v-if="conn.value !== 'Open'">Connection {{conn}}</span>

    <div v-if="selection">
      {{richnessText}} {{treeText}}
    </div>
    <div>
      <button @click="endTurn()">End Turn</button>
      <span v-if="tree && tree.Owner === you && !tree.IsDormant">
        <button v-if="tree.Size >= 1" @click="seed1(selection)">Seed</button>
        <button v-if="tree.Size < 3" @click="grow(selection)">Grow</button>
        <button v-if="tree.Size === 3" @click="sell(selection)">Sell</button>
      </span>
      <span v-if="!tree && seedSource && cell.Richness > 0">
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
    },
    richnessText() {
      if (!this.cell) {
        return ''
      }
      switch (this.cell.Richness) {
        case 0:
          return 'Unusable'
        case 2:
          return 'Soil Bonus +1'
        case 3:
          return 'Soil Bonus +2'
        default:
          return ''
      }
    },
    treeText() {
      if (!this.tree) {
        return ''
      }
      let text = ''
      if (this.tree.Owner === this.you) {
        text += 'Your '
      } else {
        text += `Opponent's `
      }
      if (this.tree.IsDormant) {
        text += 'dormant '
      }
      switch (this.tree.Size) {
        case 0:
          return text + 'seed'
        case 1:
          return text + 'sprout'
        case 2:
          return text + 'sapling'
        case 3:
          return text + 'tree'
        default:
          return text + '???'
      }
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
