<template>
  <g class="cell" @click="cellClick">
    <polygon points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87" :class="polyClass"></polygon>
    <text class="xyz"></text>
    <image v-if="href" transform="rotate(30)" :x="offsetX" :y="offsetY" :width="size" :height="size" :class="{bluePlayer: tree?.Owner, orangePlayer: !tree?.Owner}" :href="href" />
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
  </g>
</template>
<script>
export default {
  props: ['hex', 'index', 'debug', 'tree', 'cell', 'you', 'selection'],
  computed: {
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
      if (this.tree.Size === 0) {
        return '/seed.svg'
      } else if (this.tree.Size === 1) {
        return '/sprout.svg'
      } else if (this.tree.Size === 2) {
        return '/med.svg'
      } else if (this.tree.Size === 3) {
        return '/tree.svg'
      }
    },
    size() {
      if (this.href === '/tree.svg') {
        return '150'
      }
      return '80'
    },
    offsetX() {
      if (this.href === '/tree.svg') {
        return '-75'
      }
      return '-40'
    },
    offsetY() {
      if (this.href === '/tree.svg') {
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
  /*stroke-width: 3px;*/
}
.richUnusable {
  fill: gray;
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
.bluePlayer {
  filter: sepia(100%) saturate(300%) brightness(100%) hue-rotate(180deg);
}
.orangePlayer {
  filter: sepia(100%) saturate(300%) brightness(100%) hue-rotate(0deg);
}
</style>
