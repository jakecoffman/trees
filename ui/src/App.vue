<template>
  <main>
    <hex-grid/>
    <div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
  </main>
</template>
<script>
import SimpleGrid from "./components/SimpleGrid.vue";
import HexGrid from "./components/HexGrid.vue";
export default {
  components: {
    SimpleGrid,
    HexGrid
  }
}
</script>
<style>
body {
  margin: 0;
}
</style>
