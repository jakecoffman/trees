<template>
  <svg id="sun" viewBox="0 0 497 497" width="512" height="512" xmlns="http://www.w3.org/2000/svg">
    <!-- Le sun -->
    <circle fill="#ffd900" cx="248.5" cy="248.5" r="149.1"/>
    <!-- Sun shading -->
    <path fill="#fb0" d="m248.5 99.4c-3.87 0-7.705.148-11.5.437 76.972 5.871 137.6 70.187 137.6 148.663s-60.628 142.792-137.6 148.663c3.795.289 7.63.438 11.5.438 82.346 0 149.1-66.754 149.1-149.1s-66.754-149.101-149.1-149.101z"/>

    <g :transform="`rotate(${rotation})`" id="arrow">
      <!-- arrow -->
      <path fill="#ffd900" transform="translate(0 1) scale(0.3 0.3)" d="m1635.35393,777.85464l-206.85734,-187.908c-9.55306,-8.684 -19.17847,-13.28 -27.17904,-13.28c-15.46758,0 -25.03571,16.472 -25.03571,44.044l0,406.124c0,27.54 9.55607,43.98 24.98747,43.98c8.01263,0 17.4843,-4.6 27.05846,-13.308l206.9538,-187.904c13.30917,-12.104 20.6797,-28.392 20.6797,-45.884c0.00301,-17.48 -7.28311,-33.764 -20.60736,-45.864z"/>
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
