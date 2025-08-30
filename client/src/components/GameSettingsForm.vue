<script setup lang="ts">
import { type FormHTMLAttributes } from 'vue';
import type { GameSettings } from '../state';

interface Props extends /* @vue-ignore */ FormHTMLAttributes {}
defineProps<Props>();

const model = defineModel<GameSettings>({ required: true });
</script>

<template>
  <form>
    <label>
      <span>Goal</span>
      <input v-model="model.Goal" type="number" name="Goal" min="1" max="255" />
    </label>
    <label>
      <span>White Cards</span>
      <input
        v-model="model.HandWhites"
        type="number"
        name="HandWhites"
        min="1"
        max="255"
      />
    </label>
    <label>
      <span>Black Cards</span>
      <input
        v-model="model.HandBlacks"
        type="number"
        name="HandBlacks"
        min="1"
        max="255"
      />
    </label>
    <label>
      <span>Fighter Cards</span>
      <input
        v-model="model.FighterBlacks"
        type="number"
        name="FighterBlacks"
        min="1"
        max="255"
      />
    </label>
  </form>
</template>
