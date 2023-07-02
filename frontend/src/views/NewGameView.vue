<template>
  <div class="px-4 py-5 my-5 text-center">
    <h1 class="display-5 fw-bold text-body-emphasis">Let the game begin</h1>
    <p class="lead mb-4">
      Enter your name and click the button below to start playing.
    </p>
    <div class="p-4 p-md-5 border rounded-3 bg-body-tertiary col-lg-4 mx-auto">
      <div class="form-floating mb-3">
        <input
            type="text"
            class="form-control"
            id="name"
            minlength="3"
            maxlength="20"
            placeholder="Magnus Carlsen"
            :class="{
                'is-invalid': invalidName,
            }"
            v-model="name"/>
        <!-- eslint-disable-next-line -->
        <label for="name">Your name</label>
      </div>
      <button
        class="w-100 btn btn-lg btn-primary"
        @click="findMatch"
        id="findgame"
        v-html="buttonText"></button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { APIClient, EventGameStarted } from '@/api/client';

export default defineComponent({
    name: 'CreateGameView',
    data() {
        return {
            name: '',
            buttonText: 'Find a game',
            invalidName: false,
        };
    },
    methods: {
        async findMatch() {
            if (this.name.length < 3 || this.name.length > 20) {
                this.invalidName = true;
                return;
            }
            this.invalidName = false;

            try {
                this.buttonText = '<i class="fa fa-spinner fa-spin"></i> Looking for a game';
                const api = APIClient.getInstance();
                const onGameStarted = (evt: Event) => {
                    const messageEvent = (evt as MessageEvent);
                    const gameStarted: EventGameStarted = JSON.parse(messageEvent.data);

                    let playerSide = '';
                    if (gameStarted.PlayerBlack === this.name) {
                        playerSide = 'black';
                    } else if (gameStarted.PlayerWhite === this.name) {
                        playerSide = 'white';
                    }

                    if (playerSide !== '') {
                        api.forget('game:start', onGameStarted);
                        this.$router.push({
                            name: 'game',
                            state: {
                                playerName: this.name,
                                gameId: gameStarted.GameID,
                                playerSide,
                            },
                        });
                    }
                };

                await api.listen('game:start', onGameStarted);
                try {
                    await api.findMatch(this.name);
                } catch (error) {
                    console.error(error);
                    api.forget('game:start', onGameStarted);
                    this.buttonText = 'Find a game';
                }
            } catch (error) {
                console.error(error);
            }
        },
    },
});
</script>

<style>
</style>
