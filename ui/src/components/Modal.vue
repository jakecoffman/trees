<template>
  <aside class="modal">
    <slot></slot>
  </aside>
</template>
<script>
export default {}
</script>
<style scoped>
.modal {
  position: absolute;
  top: 1rem;
  left: 1rem;
  right: 1rem;
  z-index: 1;
  background: white;
  border: 1px solid black;
  border-radius: 5px;

  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 1rem;

  box-shadow: gray 5px 5px;
}
</style>
