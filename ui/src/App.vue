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
  color: #c5c5c5;
  font-size: 16pt;
}
a {
  color: #91c8ff;
}
.orange {
  color: #FFA500;
}
.dark-orange {
  color: #c48712;
}
.orange-bg {
  background: black;
}
.blue {
  color: #8fa2ff;
}
.dark-blue {
  color: #334ed9;
}
.blue-bg {
  background: black;
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
  border-radius: 5px;
}
</style>
