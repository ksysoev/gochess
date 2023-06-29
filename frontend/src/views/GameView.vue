<template>
    <div>{{opponentName}}</div>
    <div>
      <TheChessboard
        @move="onMove"
        :board-config="boardConfig"
        @board-created="(api) => (boardAPI = api)"/>
    </div>
    <div>{{playerName}}</div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from 'vue';
import {
    TheChessboard, type MoveEvent, type BoardApi,
} from 'vue3-chessboard';
import { APIClient, EventGameMove } from '@/api/client';
import 'vue3-chessboard/style.css';

export default defineComponent({
    name: 'GameView',

    data() {
        const { playerSide, playerName, gameId } = window.history.state;

        return {
            gameId,
            playerName,
            playerSide,
            opponentName: '',
            boardConfig: {
                position: 'start',
                orientation: playerSide,
                coordinates: true,
            },
        };
    },
    async created() {
        if (!window.history.state.gameId) {
            this.$router.push({ name: 'home' });
        }
    },
    components: {
        TheChessboard,
    },
});
</script>

<script setup lang="ts">
const api = APIClient.getInstance();

const gameId = window.history.state.gameId || '';
let boardAPI: BoardApi | undefined;
let LastServerMove = '';
let MyLastMove = '';

onMounted(async () => {
    const game = await api.getGame(gameId);

    boardAPI?.setPosition(game.position);
});

async function onMove(move: MoveEvent) {
    if (LastServerMove === move.san) {
        return;
    }
    MyLastMove = move.san;
    const updatedGame = await api.makeMove(gameId, move.san);
    if (boardAPI) {
        boardAPI.setPosition(updatedGame.position);
    } else {
        console.error('boardAPI is undefined');
    }
}

api.listen('game:move', (evt: Event) => {
    const messageEvent = (evt as MessageEvent);
    const moveEvent: EventGameMove = JSON.parse(messageEvent.data);
    if (MyLastMove === moveEvent.Move) {
        return;
    }
    if (moveEvent.GameID === gameId) {
        if (boardAPI) {
            LastServerMove = moveEvent.Move;
            boardAPI.move(moveEvent.Move);
        } else {
            console.error('boardAPI is undefined');
        }
    }
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
