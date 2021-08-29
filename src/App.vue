<template>
  <simple-grid/>
  <hex-grid/>
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
