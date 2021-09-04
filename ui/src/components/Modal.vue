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
  top: 40%;
  left: 1rem;
  right: 1rem;
  z-index: 1;
  background: #c9c9c9;
  color: black;
  border: 1px solid black;
  border-radius: 5px;
  margin: 1rem;
  padding: 1rem;

  display: flex;
  flex-direction: column;
  align-items: center;

  box-shadow: 10px 10px 10px rgba(0, 0, 0, 0.5);
}
</style>
