<template>
  <g class="cell">
    <text class="xyz">{{x}} {{y}} {{z}}</text>
    <polygon points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87"></polygon>
    <image
        transform="rotate(30)"
        :x="offsetX"
        :y="offsetY"
        :width="size"
        :height="size"
        style="filter: sepia(100%) saturate(300%) brightness(70%) hue-rotate(180deg);"
        :href="href" />
    <g v-if="debug">
      <text transform="rotate(-90) translate(60, 0) rotate(90) rotate(30) translate(0,10)" class="q-coord">
        {{x}}
      </text>
      <text transform="rotate(-90) translate(-30, -52) rotate(90) rotate(30) translate(0,10)" class="s-coord">
        {{y}}
      </text>
      <text transform="rotate(-90) translate(-30, 52) rotate(90) rotate(30) translate(0,10)" class="r-coord">
        {{z}}
      </text>
      <text transform="rotate(30)" style="stroke: red;">
        {{index}}
      </text>
    </g>
  </g>
</template>
<script>
export default {
  props: ['x', 'y', 'z', 'index', 'debug'],
  computed: {
    href() {
      const v = this.index % 5
      if (v === 0) {
        return 'seed.svg'
      } else if (v === 1) {
        return 'sprout.svg'
      } else if (v === 2) {
        return 'med.svg'
      } else if (v === 3) {
        return 'tree.svg'
      }
    },
    size() {
      if (this.href === 'tree.svg') {
        return '150'
      }
      return '80'
    },
    offsetX() {
      if (this.href === 'tree.svg') {
        return '-75'
      }
      return '-40'
    },
    offsetY() {
      if (this.href === 'tree.svg') {
        return '-70'
      }
      return '-40'
    }
  }
}
</script>
<style scoped>
.cell:hover {
  fill: #9b9bff;
}
</style>
