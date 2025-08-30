<script setup lang="ts">
import type { FormHTMLAttributes } from 'vue';
import RocketIcon from '../assets/rocket.svg?component';

interface Props extends /* @vue-ignore */ FormHTMLAttributes {}
defineProps<Props>();

const model = defineModel<string>();
</script>

<template>
  <form>
    <input
      v-model="model"
      type="text"
      name="name"
      placeholder="Your Name"
      required
    />
    <button type="submit">
      <RocketIcon />
      <span>Play</span>
    </button>
  </form>
</template>
