<template>
  <section class="welcome">
    <h1>Tree Game</h1>
    <p>
      Tree Game is a 2-player strategy game where you try to score
      the most points by selling mature trees before you opponent does.
    </p>
    <p>
      <router-link to="/rules">
        Read full rules
      </router-link>
    </p>
    <div class="welcome-choices">
      <button @click="$router.push('/games')">
        New Game
      </button>
      <button @click="showJoinModal()">
        Join Game
      </button>
    </div>

    <transition name="slide">
      <modal v-if="existing">
        <p>You are already in a game.</p>
        <footer class="modal-footer">
          <button @click="quit()">
            Quit
          </button>
          <button @click="$router.push(`/games/${existing}`)">
            Rejoin
          </button>
        </footer>
      </modal>
    </transition>

    <transition name="slide">
      <modal v-model="join">
        <p>Enter the ID of the game</p>
        <input type="number" v-model="code" ref="codeEntry" :disabled="code.toString().length === 6">
        <footer class="modal-footer">
          <button @click="join = false; code = ''">
            Cancel
          </button>
        </footer>
      </modal>
    </transition>
  </section>
</template>
<script>
import Modal from "../components/Modal.vue";
import {createToast} from "mosha-vue-toastify";
export default {
  components: {
    Modal
  },
  data() {
    return {
      join: false,
      code: '',
      existing: ''
    }
  },
  watch: {
    async code() {
      if (this.code.toString().length === 6) {
        const r = await fetch(`/api/rooms/${this.code}`)
        if (!r.ok) {
          createToast(`Room ${this.code} not found`, {type: 'danger', position: 'bottom-right'})
          this.code = ''
          await this.$nextTick(() => {
            this.$refs.codeEntry.focus()
          })
          return
        }
        await this.$router.push(`/games/${this.code}`)
      }
    }
  },
  async created() {
    const r = await fetch('/api/login')
    if (!r.ok) {
      return alert("Failed to login")
    }
    const me = await r.json()
    this.existing = me.Code
  },
  methods: {
    async quit() {
      const r = await fetch(`/api/rooms/${this.existing}`, {method: 'DELETE'})
      if (!r.ok) {
        return alert("Failed to quit")
      }
      await r.text()
      this.existing = ''
    },
    showJoinModal() {
      this.join = true
      this.$nextTick(() => {
        this.$refs.codeEntry.focus()
      })
    }
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
.welcome-choices {
  display: flex;
  justify-content: center;
}
</style>
