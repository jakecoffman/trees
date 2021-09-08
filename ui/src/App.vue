<template>
  <main>
    <router-view/>
  </main>
</template>
<script>
export default {}
</script>
<style>
@font-face {
  font-family: 'Roboto'; /*a name to be used later*/
  src: url('/Roboto/Roboto-Regular.ttf'); /*URL to font*/
}
h1, h2, h3, h4 {
  margin: 0;
  padding: 0;
}
html {
  font-family: Roboto,sans-serif;
  font-size: 16pt;
}
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
  background: #e1e1e1;
  border: 1px solid black;
  font-size: 16pt;
  padding: .25rem 1rem;
  border-radius: 5px;
}
p {
  margin-top: .5rem;
  margin-bottom: .5rem;
}
@media (max-height: 600px) {
  html {
    font-size: 6px;
  }
  body {
    font-size: 12pt;
  }
  button {
    font-size: 12pt;
  }
}
</style>
