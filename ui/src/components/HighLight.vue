<template>
  <transition name="highlight" mode="out-in">
    <slot/>
  </transition>
</template>
<style scoped>
.highlight-enter-active {
  animation: highlight 5s;
}
@keyframes highlight {
  0% { background-color: #31792c; }
  100% { background-color: transparent; }
}
</style>
