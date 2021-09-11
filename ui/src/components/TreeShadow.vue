<template>
  <g v-if="tree && tree.Size === showSize" :transform="`rotate(${rotation})`" class="shadow">
    <image class="shadow" v-if="tree && tree.Size > 0"
           style="opacity: 70%;"
           :transform="`rotate(30) translate(25, -50) scale(${scale}, 1)`"
           width="100"
           height="100"
           href="/cone.svg"/>
  </g>
</template>
<script>
export default {
  props: ['game', 'index', 'showSize'],
  computed: {
    rotation() {
      return this.day * -60
    },
    day() {
      return this.game.State.Day
    },
    tree() {
      return this.game.State.Trees[this.index]
    },
    scale() {
      switch (this.tree.Size) {
        case 1:
          return 2
        case 2:
          return 4
        case 3:
          return 5.6
      }
      return 0
    }
  }
}
</script>
<style scoped>
.shadow {
  pointer-events: none;
  transition: all 0.5s ease-in-out;
}
</style>
