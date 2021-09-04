<template>
  <footer>
    <span v-if="conn.value !== 'Open'">Connection {{conn}}</span>

    <div v-if="locked">Waiting for opponent's move</div>
    <div v-if="selection !== null && !seedSource">
      {{richnessText}} {{treeText}}
    </div>
    <div v-if="game.State.Day < 26" class="buttons">
      <span v-if="tree && tree.Owner === you && !tree.IsDormant && !seedSource">
        <button v-if="tree.Size >= 1" @click="seed1(selection)" :disabled="locked || seedCost > game.State.Energy[you]">
          Seed (Cost {{seedCost}})
        </button>
        <button v-if="tree.Size < 3" @click="grow(selection)" :disabled="locked || growthCost > game.State.Energy[you]">
          Grow (Cost {{growthCost}})
        </button>
        <button v-if="tree.Size === 3" @click="sell(selection)" :disabled="locked || growthCost > game.State.Energy[you]">
          Sell (Cost {{growthCost}})
        </button>
      </span>
      <span v-if="!tree && seedSource && cell?.Richness > 0">
        <button @click="seed2(selection)">Seed Here</button>
      </span>
      <div v-if="seedSource" class="flex column center">
        <p v-if="tree || cell?.Richness === 0">Select a location to seed</p>
        <button @click="seedSource = null">Cancel Seed</button>
      </div>
      <button @click="endTurn()" v-if="!seedSource" :disabled="locked">End Turn</button>
    </div>
    <span v-else>
      <span v-if="game.State.Score[0] === game.State.Score[1]">Tie Game!</span>
      <span v-else-if="game.State.Score[0] > game.State.Score[1]">Orange Wins!</span>
      <span v-else>Blue Wins!</span>
    </span>
  </footer>
</template>
<script>
import {growthCost, seedCost} from "../assets/cost";

export default {
  props: ['game', 'selection', 'you'],
  inject: ['ws', 'conn'],
  data() {
    return {
      seedSource: null,
      locked: false
    }
  },
  computed: {
    cell() {
      if (this.selection === null) {
        return null
      }
      return this.game.State.Board.Cells[this.selection]
    },
    tree() {
      if (this.selection === null) {
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
    },
    growthCost() {
      if (!this.tree) {
        return 0
      }
      return growthCost(this.you, this.game.State, this.tree)
    },
    seedCost() {
      return seedCost(this.you, this.game.State)
    }
  },
  methods: {
    endTurn() {
      this.locked = true
      this.ws.value.send(JSON.stringify({Kind: 'end'}))
    },
    seed1(selection) {
      this.seedSource = selection
    },
    seed2(selection) {
      this.locked = true
      this.ws.value.send(JSON.stringify({Kind: 'seed', Source: this.seedSource, Target: selection}))
      this.seedSource = null
    },
    sell(selection) {
      this.locked = true
      this.ws.value.send(JSON.stringify({Kind: 'sell', Target: selection}))
    },
    grow(selection) {
      this.locked = true
      this.ws.value.send(JSON.stringify({Kind: 'grow', Target: selection}))
    },
    unlock() {
      this.locked = false
    }
  }
}
</script>
<style scoped>
.buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  justify-content: space-between;
  align-items: center;
}
footer {
  margin-top: 1rem;

  display: flex;
  gap: .5rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.flex {
  display: flex;
}
.column {
  flex-direction: column;
}
.center {
  align-items: center;
}
</style>
