<template>
    <div class="card ">
        <div class="card-body">
            {{ opponentName }}
        </div>
    </div>
    <div>
      <TheChessboard
        @move="onMove"
        :board-config="boardConfig"
        @board-created="(api) => (boardAPI = api)"/>
    </div>
    <div class="card">
        <div class="card-body">
            {{ playerName }}
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
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
            boardConfig: {
                position: 'start',
                orientation: playerSide,
                coordinates: true,
                movable: {
                    free: false,
                    color: playerSide,
                },
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
const { playerSide, gameId } = window.history.state;
const opponentName = ref('');

let boardAPI: BoardApi | undefined;
let LastServerMove = '';
let MyLastMove = '';

onMounted(async () => {
    const game = await api.getGame(gameId);
    boardAPI?.setPosition(game.position);
    if (playerSide === 'white') {
        opponentName.value = game.playerBlack;
    } else {
        opponentName.value = game.playerWhite;
    }
});

async function onMove(move: MoveEvent) {
    if (!boardAPI) {
        return;
    }

    if (LastServerMove === move.san) {
        return;
    }

    if (move.color !== playerSide.toLowerCase().charAt(0)) {
        boardAPI?.setPosition(move.before);
        return;
    }

    MyLastMove = move.san;
    try {
        const updatedGame = await api.makeMove(gameId, move.san);
        boardAPI.setPosition(updatedGame.position);
    } catch (error) {
        boardAPI?.setPosition(move.before);
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
}
</style>
