<template>
    <div>{{black}}</div>
    <div>
      <TheChessboard @move="onMove" :board-config="boardConfig" />
    </div>
    <div>{{white}}</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { TheChessboard, type MoveEvent, type BoardConfig } from 'vue3-chessboard';
import 'vue3-chessboard/style.css';

export default defineComponent({
  name: 'GameView',
  props: {
    white: String,
    black: String,
    gameId: String,
  },
  components: {
    TheChessboard,
  },
});
</script>

<script setup lang="ts">
let gameId: string | null = null;
const boardConfig: BoardConfig = {};

fetch('http://localhost:8081/game', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
})
  .then((response) => response.json())
  .then((data) => {
    gameId = data.id;
  })
  .catch((error) => {
    console.error(error);
  });

function onMove(move: MoveEvent) {
  console.log(move);
  if (gameId) {
    fetch(`http://localhost:8081/game/${gameId}/move`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ move: move.san }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.error(error);
      });
  }
}
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
