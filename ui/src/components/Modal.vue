<template>
  <div>
    <div v-if="open || modelValue" class="smokescreen" @click="close"></div>
    <transition name="slide">
      <aside v-if="open || modelValue" class="modal">
        <slot></slot>
      </aside>
    </transition>
  </div>
</template>
<script>
export default {
  // use v-model if you have a variable, otherwise use open.
  // if you use both weird things will happens
  props: ['modelValue', 'open'],
  methods: {
    close() {
      if (this.open !== null) {
        this.$emit('update:modelValue', false)
      }
    }
  }
}
</script>
<style scoped>
.modal {
  position: absolute;
  top: 50%;
  transform: translate(0, -50%);
  left: 1rem;
  right: 1rem;
  z-index: 2;
  background: #c9c9c9;
  color: black;
  border: 1px solid black;
  border-radius: 5px;
  padding: 1rem;

  display: flex;
  flex-direction: column;
  align-items: center;

  box-shadow: 10px 10px 10px rgba(0, 0, 0, 0.5);
}
.smokescreen {
  position: fixed;
  z-index: 1;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  opacity: 50%;
  background: black;
}
</style>
