<template>
  <div>
    <TheChessboard />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { TheChessboard } from 'vue3-chessboard';
import 'vue3-chessboard/style.css';

export default defineComponent({
  name: 'App',
  components: {
    TheChessboard,
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
