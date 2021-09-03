<template>
  <main>
    <router-view/>
  </main>
</template>
<script>
export default {}
</script>
<style>
body {
  margin: 0;
  background: black;
  color: white;
  font-size: 16pt;
}
.orange {
  color: orange;
}
.orange-bg {
  background: orange;
}
.blue {
  color: #8fa2ff;
}
.blue-bg {
  background: #8fa2ff;
}
.flex {
  display: flex;
}
.grow {
  flex-grow: 1;
}
button {
  background: white;
  border: 1px solid black;
  font-size: 16pt;
  padding: .25rem 1rem;
}
</style>
