<script setup lang="ts">
import RocketIcon from '../assets/rocket.svg?component';

const model = defineModel<string>({ required: true });
</script>

<template>
  <form>
    <input v-model="model" type="text" placeholder="Your Name" required />
    <button type="submit">
      <RocketIcon />
      <span>Play</span>
    </button>
  </form>
</template>
