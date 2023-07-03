<template>
    <div class="card">
        <div class="card-body text-start">
            <strong>{{ opponentName }}</strong>
        </div>
    </div>
    <div>
    <TheChessboard
        @move="onMove"
        @checkmate="handleCheckmate"
        :board-config="boardConfig"
        @board-created="(api) => (boardAPI = api)"/>
    </div>
    <div class="card">
        <div class="card-body text-left text-start font-weight-bolder">
            <strong>{{ playerName }}</strong>
        </div>
    </div>

    <!-- Modal -->
    <div
        class="modal fade"
        id="exampleModal"
        tabindex="-1"
        role="dialog"
        aria-labelledby="exampleModalLabel"
        aria-hidden="true">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
        <div class="modal-header">
            <h5 class="modal-title" id="exampleModalLabel">Game Result</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
            </button>
        </div>
        <div class="modal-body">
            {{ opponentName }}
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        </div>
        </div>
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
const gameResult = ref('');

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

    handleState(game.State);
});

async function handleCheckmate(isMated) {
    if (!boardAPI) {
        return;
    }
    onMove(boardAPI.getLastMove());
}

async function onMove(move: MoveEvent) {
    if (!move || !boardAPI) {
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
        if (updatedGame.position !== move.after) {
            boardAPI?.setPosition(updatedGame.position);
        }
    } catch (error) {
        boardAPI?.setPosition(move.before);
    }
}

api.listen('game:move', (evt: Event) => {
    const messageEvent = (evt as MessageEvent);
    const moveEvent: EventGameMove = JSON.parse(messageEvent.data);

    if (moveEvent.GameID !== gameId) {
        return;
    }

    if (MyLastMove !== moveEvent.Move) {
        if (boardAPI) {
            LastServerMove = moveEvent.Move;
            boardAPI.move(moveEvent.Move);
        } else {
            console.error('boardAPI is undefined');
        }
    }

    handleState(moveEvent.State);
});

function handleState(state) {
    switch (state) {
    case 'white_won':
        gameResult.value = 'White won';
        break;
    case 'black_won':
        gameResult.value = 'Black won';
        break;

    case 'draw':
        gameResult.value = 'Draw';
        break;
    default:
        break;
    }

    // TODO Replace it with a modal
    if (gameResult.value) {
        alert(gameResult.value);
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
}
</style>
