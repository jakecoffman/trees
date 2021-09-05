<template>
  <g class="cell" @click="cellClick">
    <polygon class="hex" points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87" :class="polyClass"></polygon>
    <image v-if="href" transform="rotate(30)" :x="offsetX" :y="offsetY" :width="size" :height="size" :href="href" />
    <text v-if="tree?.IsDormant" transform="rotate(30) translate(-50, 30)" style="font-size: 72pt">💤</text>
    <g v-if="debug">
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
  props: ['hex', 'index', 'debug', 'game', 'you', 'selection'],
  computed: {
    tree() {
      return this.game.State.Trees[this.index]
    },
    cell() {
      return this.game.State.Board.Cells[this.index]
    },
    shadow() {
      if (this.game.State.Shadows[3] == null) {
        // not sure why but sometimes shadows is [null, null, null, null]
        return 0
      }
      if (this.game.State.Shadows[3].includes(this.index)) {
        return 3
      }
      if (this.game.State.Shadows[2].includes(this.index)) {
        return 2
      }
      if (this.game.State.Shadows[1].includes(this.index)) {
        return 1
      }
      return 0
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
    href() {
      if (!this.tree) {
        return null
      }
      let path = "/"
      if (this.tree.Owner) {
        path += 'blue'
      } else {
        path += 'orange'
      }
      if (this.tree.Size === 0) {
        path += '/seed.svg'
      } else if (this.tree.Size === 1) {
        path += '/sprout.svg'
      } else if (this.tree.Size === 2) {
        path += '/med.svg'
      } else if (this.tree.Size === 3) {
        path += '/tree.svg'
      }
      return path
    },
    size() {
      if (this.tree.Size === 3) {
        return '150'
      }
      return '80'
    },
    offsetX() {
      if (this.tree.Size === 3) {
        return '-75'
      }
      return '-40'
    },
    offsetY() {
      if (this.tree.Size === 3) {
        return '-70'
      }
      return '-40'
    }
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
