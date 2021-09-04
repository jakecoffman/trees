<template>
  <svg id="sun" height="512" viewBox="0 0 497 497" width="512" xmlns="http://www.w3.org/2000/svg">
    <!-- Le sun -->
    <circle fill="#ffd900" cx="248.5" cy="248.5" r="149.1"/>
    <!-- Sun shading -->
    <path fill="#fb0" d="m248.5 99.4c-3.87 0-7.705.148-11.5.437 76.972 5.871 137.6 70.187 137.6 148.663s-60.628 142.792-137.6 148.663c3.795.289 7.63.438 11.5.438 82.346 0 149.1-66.754 149.1-149.1s-66.754-149.101-149.1-149.101z"/>

    <g :transform="`rotate(${rotation})`" id="arrow">
      <!-- ray -->
      <path fill="#ffd900" d="m447.3 231.933c-9.149 0-16.567 7.417-16.567 16.567s7.417 16.567 16.567 16.567h33.133c9.15 0 16.567-7.417 16.567-16.567s-7.417-16.567-16.567-16.567z"/>
      <!-- ray shading -->
      <path fill="#fb0" d="m480.433 231.933h-23c9.15 0 16.567 7.417 16.567 16.567 0 9.149-7.417 16.567-16.567 16.567h23c9.15 0 16.567-7.417 16.567-16.567 0-9.149-7.417-16.567-16.567-16.567z"/>
    </g>
  </svg>
</template>
<script>
export default {
  props: ['day'],
  data() {
    return {
      rotation: this.day * -60
    }
  },
  watch: {
    day() {
      this.rotation = this.day * -60
    }
  },
  computed: {
    degrees() {
      const r = this.day %6
      switch (r) {
        case 0:
          return 0
        case 1:
          return -30
        case 2:
          return -130
        case 3:
          return -180
        case 4:
          return -220
        case 5:
          return -310
        default:
          return 0
      }
    }
  }
}
</script>
<style scoped>
#sun {
  pointer-events: none;
}
#arrow {
  transform-origin: center;
  transition: all 0.5s ease-in-out;
}
</style>
