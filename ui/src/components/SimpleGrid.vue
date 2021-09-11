<template>
  <div class="board">
    <div class="row" v-for="row of board" :key="row">
      <span v-for="(item, index) of row" :key="index">{{item}}</span>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      board: [
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
          [0, 1, 2, 3, 4, 5, 6],
      ]
    }
  }
}
</script>
<style scoped>
.board {
  display: grid;
  grid-template-rows: repeat(7, 32px);
}
.row {
  display: grid;
  grid-template-columns: repeat(7, 32px);
}
span {
  border: 1px solid black;
}
</style>
