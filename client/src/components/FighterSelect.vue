<script setup lang="ts">
import { type FormHTMLAttributes } from 'vue';
import type { Fighter } from '../state';

interface Props extends /* @vue-ignore */ FormHTMLAttributes {
  disabled?: boolean;
  fighters: Fighter[];
}
defineProps<Props>();

const model = defineModel<number>();
</script>

<template>
  <form>
    <label
      v-for="({ Black, White, Tiebreak, Player }, i) in fighters"
      :key="Player"
      class="fighter"
    >
      <input
        v-model.number="model"
        type="radio"
        name="fighter"
        :value="i + 1"
        :disabled="disabled"
      />
      <span class="card white">{{ White }}</span>
      <span v-for="text in Black" :key="text" class="card black">
        {{ text }}
      </span>
      <span v-if="!!Tiebreak" class="card white">{{ Tiebreak }}</span>
    </label>
  </form>
</template>
