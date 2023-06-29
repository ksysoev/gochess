<template>
  <div class="px-4 py-5 my-5 text-center">
    <h1 class="display-5 fw-bold text-body-emphasis">Welcome to Go Chess</h1>
    <div class="col-lg-6 mx-auto">
      <p class="lead mb-4">
        Get ready to test your strategic skills and start playing now.
        Click the button below to begin your chess journey!
      </p>
      <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
        <router-link to="/new-game" class="btn btn-primary btn-lg px-4 gap-3">
          Start Game
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'HomePage',
});
</script>

<style>
</style>
