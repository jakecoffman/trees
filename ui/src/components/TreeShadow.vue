<template>
  <g :transform="`rotate(${rotation})`" class="shadow">
    <image class="shadow" v-if="tree && tree.Size > 0"
           style="opacity: 50%;"
           :transform="`rotate(30) translate(0, -50) scale(${scale}, 1)`"
           width="100"
           height="100"
           href="/cone.svg"/>
  </g>
</template>
<script>
export default {
  props: ['game', 'index'],
  data() {
    return {
      rotation: this.day * -60
    }
  },
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
          return 6
      }
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
