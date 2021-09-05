<template>
  <g class="tree" v-if="tree && tree.Size === showSize">
    <image v-if="href" transform="rotate(30)" :x="offsetX" :y="offsetY" :width="size" :height="size" :href="href" />
    <text v-if="tree?.IsDormant" transform="rotate(30) translate(-50, 30)" style="font-size: 72pt">💤</text>
  </g>
</template>
<script>
export default {
  props: ['game', 'index', 'showSize'],
  computed: {
    tree() {
      return this.game.State.Trees[this.index]
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
  }
}
</script>
<style scoped>
.tree {
  pointer-events: none;
}
</style>
