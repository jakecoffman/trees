<template>
  <g class="q-axis-same">
    <text class="xyz">{{x}} {{y}} {{z}}</text>
    <polygon points="100,0 50,-87 -50,-87 -100,-0 -50,87 50,87"></polygon>
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
  props: ['x', 'y', 'z', 'index']
}
</script>
<style scoped>
.q-axis-same {
  stroke: blue;
}
</style>
