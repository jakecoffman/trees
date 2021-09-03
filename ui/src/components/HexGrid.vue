<template>
  <section class="game" :class="{'blue-bg': !you, 'orange-bg': you}">
    <svg class="board" viewBox="-610 -560 1220 1115">
      <g transform="rotate(-30)" fill="white" stroke="black">
        <g class="grid" id="grid">
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <hex-cell
                :hex="hex"
                :index="index"
                :tree="game.State.Trees[index]"
                :cell="game.State.Board.Cells[index]"
                :you="you"
                :selection="selection"
                @select="select"
            />
          </g>
        </g>
      </g>
    </svg>
    <div style="color: black;">Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
  </section>
</template>
<script>
import HexCell from "./HexCell.vue";
import {grid} from "../assets/grid";

export default {
  components: {
    HexCell
  },
  props: ['game', 'you', 'selection'],
  data() {
    return {
      grid
    }
  },
  methods: {
    select(index) {
      this.$emit('select', index)
    }
  }
}

</script>
<style scoped>
.game {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.board {
  max-width: 600px;
  max-height: 90vh;
}
</style>
