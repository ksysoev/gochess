<template>
  <div>
    <div v-if="!showGameView">
      <label for="username">
        <input
          placeholder="Enter your username"
          type="text"
          id="username"
          v-model="username"/>
      </label>
      <button @click="findMatch">Find a match</button>
    </div>
    <GameView v-if="showGameView" :white="white" :black="black" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import GameView from './components/GameView.vue';
import { APIConfig } from './api/config.js';

export default defineComponent({
  name: 'App',
  components: {
    GameView,
  },
  data() {
    return {
      username: '',
      white: '',
      black: '',
      showGameView: false,
    };
  },
  methods: {
    async findMatch() {
      try {
        const response = await fetch( APIConfig.baseURL + '/match', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ username: this.username }),
        });
        const data = await response.json();

        console.log(data);
        this.white = data.white;
        this.black = data.black;
        this.showGameView = true;
      } catch (error) {
        console.error(error);
      }
    },
  },
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
