<template>
  <g class="cell" @click="cellClick">
    <polygon class="hex" points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87" :class="polyClass"></polygon>
    <g v-if="1===2">
      <text transform="rotate(-90) translate(60, 0) rotate(90) rotate(30) translate(0,10)" class="q-coord">
        {{hex.x}}
      </text>
      <text transform="rotate(-90) translate(-30, -52) rotate(90) rotate(30) translate(0,10)" class="s-coord">
        {{hex.y}}
      </text>
      <text transform="rotate(-90) translate(-30, 52) rotate(90) rotate(30) translate(0,10)" class="r-coord">
        {{hex.z}}
      </text>
      <text transform="rotate(30)" style="stroke: red;">
        {{index}}
      </text>
    </g>
<!--    <text>{{hex.x}} {{hex.y}} {{hex.z}}</text>-->
  </g>
</template>
<script>
export default {
  props: ['hex', 'index', 'game', 'selection'],
  computed: {
    cell() {
      return this.game.State.Board.Cells[this.index]
    },
    polyClass() {
      if (this.selection === this.index) {
        return {selected: true}
      }
      const classes = {}

      switch (this.cell.Richness) {
        case 0:
          classes.richUnusable = true
          break
        case 1:
          classes.richLow = true
          break
        case 2:
          classes.richMed = true
          break
        case 3:
          classes.richHigh = true
          break
      }
      return classes
    },
  },
  methods: {
    cellClick() {
      this.$emit('select', this.index)
    }
  }
}
</script>
<style scoped>
.selected {
  opacity: 50%;
  fill: red;
}
.richUnusable {
  fill: #ffdd8a;
}
.richLow {
  fill: darkseagreen;
}
.richMed {
  fill: lightgreen;
}
.richHigh {
  fill: green;
}
.hex {
  transition: .5s filter linear;
}
</style>
