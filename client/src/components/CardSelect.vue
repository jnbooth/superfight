<script setup lang="ts">
import { type FieldsetHTMLAttributes } from 'vue';

interface Props extends /* @vue-ignore */ FieldsetHTMLAttributes {
  cards: string[];
  name: string;
}
defineProps<Props>();

const model = defineModel<number>();
</script>

<template>
  <fieldset>
    <label v-for="(text, i) in cards" :key="text">
      <input v-model.number="model" type="radio" :name="name" :value="i + 1" />
      <span :class="`card ${name}`">{{ text }}</span>
    </label>
  </fieldset>
</template>
