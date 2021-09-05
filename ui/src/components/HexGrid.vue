<template>
  <section class="game" :class="{'blue-bg': you, 'orange-bg': !you}">
    <the-sun class="sun" :day="game.State.Day"/>
    <svg class="board" viewBox="-610 -560 1220 1115">
      <g transform="rotate(-30)" fill="white" stroke="black">
        <g class="grid" id="grid">
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <hex-cell :hex="hex" :index="index" :game="game" :selection="selection" @select="select"/>
          </g>
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-image :index="index" :game="game" :show-size="0"/>
          </g>
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-image :index="index" :game="game" :show-size="1"/>
          </g>
          <g class="shadow" v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-shadow :index="index" :game="game" :show-size="1"/>
          </g>
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-image :index="index" :game="game" :show-size="2"/>
          </g>
          <g class="shadow" v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-shadow :index="index" :game="game" :show-size="2"/>
          </g>
          <g v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-image :index="index" :game="game" :show-size="3"/>
          </g>
          <g class="shadow" v-for="(hex, index) of grid" :key="index" :transform="`translate(${hex.tX},${hex.tY})`">
            <tree-shadow :index="index" :game="game" :show-size="3"/>
          </g>
        </g>
      </g>
    </svg>
  </section>
</template>
<script>
import HexCell from "./HexCell.vue";
import {grid} from "../assets/grid";
import TheSun from "./TheSun.vue";
import TreeShadow from "./TreeShadow.vue";
import TreeImage from "./TreeImage.vue";

export default {
  components: {
    TreeImage,
    TreeShadow,
    TheSun,
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
  max-height: 70vh;
}

.sun {
  position: absolute;
  right: 0;

  width: 100px;
  height: 100px;
}

@media (max-width: 600px) {
  .sun {
    width: 75px;
    height: 75px;
  }
}
</style>
