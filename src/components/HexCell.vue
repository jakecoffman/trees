<template>
  <g class="cell">
    <text class="xyz">{{x}} {{y}} {{z}}</text>
    <polygon points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87"></polygon>
    <image
        transform="rotate(30)"
        x="-35"
        y="-50"
        width="80"
        height="80"
        :href="href" />
    <g>
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
  props: ['x', 'y', 'z', 'index'],
  computed: {
    href() {
      const v = this.index % 4
      if (v === 0) {
        return 'seed.svg'
      } else if (v === 1) {
        return 'sprout.svg'
      } else if (v === 2) {
        return 'med.svg'
      } else {
        return 'tree.svg'
      }
    }
  }
}
</script>
<style scoped>
.cell:hover {
  fill: #9b9bff;
}
</style>
