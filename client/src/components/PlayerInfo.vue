<script setup lang="ts">
import type { LiHTMLAttributes } from 'vue';
import ClockIcon from '../assets/clock.svg?component';
import StarIcon from '../assets/star.svg?component';

interface Props extends /* @vue-ignore */ LiHTMLAttributes {
  waiting: boolean;
  points: number;
  you: boolean;
}
defineProps<Props>();
</script>

<template>
  <li :class="you ? 'you' : undefined">
    <span class="state">
      <ClockIcon v-if="waiting" />
    </span>
    <span class="name"><slot /></span>
    <span class="score">
      {{ points }}
      <StarIcon />
    </span>
  </li>
</template>
