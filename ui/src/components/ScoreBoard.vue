<template>
  <div :class="{orange: you === 0, blue: you === 1}">
    {{ you === 0 ? 'You' : 'Opponent' }}
    <high-light>
      <span :key="score">Score: {{score}}</span>
    </high-light>
    <high-light>
      <span :key="energy">
        ☀️{{ energy }}
      </span>
    </high-light>
  </div>
</template>
<script>
import HighLight from "./HighLight.vue"
export default {
  components: {HighLight},
  props: ['you', 'score', 'energy']
}
</script>
<style scoped>
.highlight-enter-active {
  animation: highlight 5s;
}
@keyframes highlight {
  0% { background-color: #a19c87; }
  100% { background-color: transparent; }
}
</style>
