<template>
  <section class="welcome">
    <h1>Tree Game</h1>
    <p>
      Tree Game is a 2-player strategy game where you try to score
      the most points by selling mature trees before you opponent does.
    </p>
    <ul>
      <li>
        <router-link to="/rules">Read full rules</router-link>
      </li>
      <li>
        <router-link to="/game">Create New Multiplayer Game</router-link>
      </li>
      <li>
        <button @click="join = true">Join Game</button>
      </li>
      <li>
        Create New Game Against Bot (Coming Soon)
      </li>
      <li>
        Log In (Coming Soon)
      </li>
    </ul>

    <transition name="slide">
      <modal v-if="join">
        <p>Enter the ID of the game</p>
        <input type="number" v-model="code">
        <footer class="modal-footer">
          <button @click="join = false; code = ''">Cancel</button>
          <button @click="$router.push(`/game/${code}`)" :disabled="code.toString().length !== 6">Join</button>
        </footer>
      </modal>
    </transition>
  </section>
</template>
<script>
import Modal from "../components/Modal.vue";
export default {
  components: {
    Modal
  },
  data() {
    return {
      join: false,
      code: ''
    }
  },
  created() {
    fetch('/api/login')
  }
}
</script>
<style scoped>
.welcome {
  margin: 1rem;
}
.modal-footer {
  margin-top: 1rem;
  display: flex;
  gap: 1rem;
}

</style>
