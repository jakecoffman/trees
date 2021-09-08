<template>
  <article>
    <game-header :game="game" :you="you"/>
    <hex-grid :game="game" :you="0" :selection="selection" @select="selection = $event"/>
    <p v-if="game.State.Day === 0">
      Welcome to the tutorial! You are the orange player and you have 3 ☀️ Energy. Press "End Turn" to
      wait until the next day. Your sprout will earn you 1 ☀️ Energy.
    </p>
    <div v-if="game.State.Day === 1">
      <p v-if="game.State.Energy[0] === 3">
        Great, now you have 3 ☀️ Energy. Click on your sprout and grow it.
      </p>
      <p v-else>
        Your sprout became a sapling and will now get 2 ☀️ Energy each morning. It is also dormant 💤 and can't
        be used again today. Press End Turn to go the next day.
      </p>
    </div>
    <div v-if="game.State.Day === 2">
      <p v-if="seedSource === null && !game.State.Trees[0].IsDormant">
        Now it's time to seed. Select your tree and pick Seed, which is free to do since you have no other seeds.
      </p>
      <p v-if="seedSource !== null && !game.State.Trees[0].IsDormant">
        Your sapling can throw a seed 2 spaces. Try to throw the seed where the sapling will never shade it.
        Select the hex where you want it to go and press "Seed Here".
      </p>
      <p v-if="game.State.Trees[0].IsDormant">
        Great job! Seeds don't earn ☀️ Energy but you can grow them into bigger trees that do. End Turn to continue.
      </p>
    </div>
    <p v-if="game.State.Day === 4">
      Notice how the trees are casting shadows? If your tree is shaded, then your tree won't gather ☀️ Energy
      that day.
    </p>
    <p v-else-if="!soldLargeTree && game.State.Day > 2 && game.State.Trees.filter(t => t && t.Size === 3).length === 0">
      Continue growing your trees until you have a large tree.
    </p>
    <p v-if="game.State.Day > 4 && game.State.Trees.filter(t => t && t.Size === 3).length === 1">
      To win the game you must have more points than your opponent on day 24. To get points you sell Large trees.
      You will gain points equal to the Nutrients. Sell your Large tree.
    </p>
    <div v-if="soldLargeTree">
      <p>
        You gained 20 points and the nutrient pool decreased by 1.
        If you wait too long to sell, the nutrients will all be gone.
        The game ends on Day 24.
      </p>
      <p>
        You completed the tutorial!
      </p>
      <div class="tutorial-controls">
        <button @click="$router.push('/')">Go Home</button>
      </div>
    </div>
    <div class="tutorial-controls" v-if="!soldLargeTree">
      <div class="tutorial-controls" v-if="tree && tree.Owner === you && !tree.IsDormant && seedSource === null">
        <button v-if="tree.Size >= 1" @click="seed1(selection)" :disabled="seedCost > game.State.Energy[you]">
          Seed (Cost {{seedCost}})
        </button>
        <button v-if="tree.Size < 3" @click="grow(selection)" :disabled="growthCost > game.State.Energy[you]">
          Grow (Cost {{growthCost}})
        </button>
        <button v-if="tree.Size === 3" @click="sell(selection)" :disabled="growthCost > game.State.Energy[you]">
          Sell (Cost {{growthCost}})
        </button>
      </div>
      <span v-if="!tree && seedSource !== null">
        <button @click="seed2(selection)">Seed Here</button>
      </span>
      <button @click="endTurn()" v-if="seedSource === null" :disabled="game.State.Day === 1 && game.State.Energy[0] === 3">
        End Turn
      </button>
    </div>
  </article>
</template>
<script>
import HexGrid from "../components/HexGrid.vue";
import GameHeader from "../components/GameHeader.vue";
import {growthCost, seedCost} from "../assets/cost";
import { createToast } from 'mosha-vue-toastify';

const cells = {}
for (let i = 0; i < 37; i++) {
  if (i < 7) {
    cells[i] = {Richness: 2}
  } else if (i < 19) {
    cells[i] = {Richness: 1}
  } else {
    cells[i] = {Richness: 0}
  }
}

export default {
  components: {
    GameHeader,
    HexGrid,
  },
  data() {
    return {
      soldLargeTree: false,
      seedSource: null,
      you: 0,
      selection: null,
      game: {
        Players: [{Connected: true},{Connected: true}],
        Code: 'Tutorial',
        State: {
          Day: 0,
          Nutrients: 20,
          Score: [0,0],
          Energy: [2,0],
          Board: {
            Cells: cells
          },
          Trees: [{
            Size: 1,
            Owner: 0
          }]
        }
      }
    }
  },
  computed: {
    tree() {
      if (this.selection === null) {
        return null
      }
      return this.game.State.Trees[this.selection.toString()]
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
    seed1(target) {
      this.seedSource = target
    },
    seed2(target) {
      const tree = this.game.State.Trees[target]
      if (tree) {
        return
      }
      const sourceTree = this.game.State.Trees[this.seedSource]
      if (sourceTree.IsDormant) {
        return
      }
      sourceTree.IsDormant = true
      this.game.State.Energy[0] -= seedCost(0, this.game.State)
      this.game.State.Trees[target] = {
        Size: 0,
        Owner: 0,
        IsDormant: true
      }
      this.seedSource = null
    },
    grow(target) {
      const tree = this.game.State.Trees[target]
      this.game.State.Energy[0] -= growthCost(0, this.game.State, tree)
      tree.Size++
      tree.IsDormant = true
    },
    sell(target) {
      const tree = this.game.State.Trees[target]
      if (!tree) {
        return
      }
      if (tree.IsDormant) {
        return
      }
      this.game.State.Energy[0] -= growthCost(0, this.game.State, tree)
      this.game.State.Score[0] += this.game.State.Nutrients
      this.game.State.Nutrients--
      this.game.State.Trees[target] = null
      this.soldLargeTree = true
      createToast('Tutorial Complete!')
    },
    endTurn() {
      this.game.State.Day++
      // gather sun
      for (const tree of this.game.State.Trees) {
        if (!tree) {
          continue
        }
        tree.IsDormant = false
        if (tree.Owner === this.you) {
          this.game.State.Energy[0] += tree.Size
        } else {
          this.game.State.Energy[1] += tree.Size
        }
      }
    }
  }
}
</script>
<style scoped>
.tutorial-controls {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
}
</style>
