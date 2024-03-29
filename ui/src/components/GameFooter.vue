<template>
  <footer>
    <img class="settings" src="/gear.svg" @click="openSettings = !openSettings">

    <div v-if="locked && game.State.Day < 24">Waiting for opponent's move</div>
    <div v-if="selection === null && game.State.Day < 24">
      Select the <span v-if="you" class="blue">blue</span><span v-else class="orange">orange</span> plants to play!
    </div>
    <div v-if="selection !== null && seedSource === null">
      {{richnessText}} {{treeText}}
    </div>
    <div v-if="game.State.Day < 24" class="buttons">
      <span v-if="tree && tree.Owner === you && !tree.IsDormant && seedSource === null">
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
      <span v-if="!tree && seedSource !== null && cell?.Richness > 0">
        <button @click="seed2(selection)">Seed Here</button>
      </span>
      <div v-if="seedSource !== null" class="flex column center">
        <p v-if="tree || cell?.Richness === 0">Select a location to seed</p>
        <button @click="seedSource = null">Cancel Seed</button>
      </div>
      <button @click="endTurn()" v-if="seedSource === null" :disabled="locked">End Turn</button>
    </div>
    <span v-else>
      <span v-if="game.State.Score[0] === game.State.Score[1]">Tie Game!</span>
      <span v-else-if="game.State.Score[0] > game.State.Score[1]">Orange Wins!</span>
      <span v-else>Blue Wins!</span>
    </span>

    <modal v-model="openSettings" class="settingsModal">
      <div class="settingsModal">
        <h3>Settings</h3>
        <button @click="$router.push(`/rules`)">
          View Rules
        </button>
        <button @click="quit()">
          Quit Game
        </button>
        <button @click="openSettings = false">
          Close Settings
        </button>
        <div style="color: gray; font-size: 10pt;">
          Created by Jake Coffman<br/>
          Icons made by Freepik from www.flaticon.com
        </div>
      </div>
    </modal>
  </footer>
</template>
<script>
import {growthCost, seedCost} from "../assets/cost";
import Modal from "./Modal.vue";

export default {
  components: {Modal},
  props: ['game', 'selection', 'you'],
  inject: ['ws', 'conn'],
  data() {
    return {
      seedSource: null,
      locked: false,
      openSettings: false
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
    },
    async quit() {
      await fetch(`/api/rooms/${this.game.Room}`, {credentials: 'include', method: 'DELETE'})
      await this.$router.push('/')
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
.just-center {
  justify-content: center;
}
.settings {
  position: fixed;
  right: 1rem;
  bottom: 1rem;
  width: 25px;
  height: 25px;
  cursor: pointer;
  z-index: 9001;
}
.settingsModal {
  z-index: 50;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
