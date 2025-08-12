<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { defaultGameState } from './state';

const etag = ref('');
const gamestate = ref(defaultGameState);
const playerName = ref('');
const playerIndex = ref(-1);
const optFighter = ref('0');
const optWhite = ref('');
const optBlack = ref('');

async function callApi(
  method: string,
  endpoint: string,
  body?: Record<string, string>,
): Promise<void> {
  const res = await fetch(endpoint, {
    method,
    body: body ? new URLSearchParams(body) : undefined,
    headers: [['If-None-Match', etag.value]],
  });
  if (res.status === 304) {
    return;
  }
  gamestate.value = await res.json();
  etag.value = res.headers.get('ETag') ?? etag.value;
  if (playerIndex.value !== -1) {
    optFighter.value =
      gamestate.value.Players[playerIndex.value].Vote.toString();
  }
}

async function poll(): Promise<void> {
  await callApi('GET', '/api/poll');
}

async function join(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/join', {
    name: playerName.value,
  });
  playerIndex.value = gamestate.value.Players.findIndex(
    player => player.Name === playerName.value,
  );
}

async function choose(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/choose', {
    player: playerIndex.value.toString(),
    white: optWhite.value,
    black: optBlack.value,
  });
}

async function vote(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/vote', {
    player: playerIndex.value.toString(),
    fighter: optFighter.value,
  });
}

/*
async function resetGame(): Promise<void> {
  await callApi('POST', '/api/reset');
}
*/

function useIntervalFn(cb: () => void, ms: number) {
  let int: number = 0;
  onMounted(() => (int = setInterval(cb, ms)));
  onUnmounted(() => clearInterval(int));
}

useIntervalFn(poll, 1000);

poll();
</script>

<template>
  <div id="players">
    <div
      v-for="{ Active, Name, Points, Vote, White } in gamestate.Players"
      :key="Name"
      :class="Active ? 'player active' : 'player'"
    >
      <span>{{ Name }}</span>
      <span>{{ Points }}</span>
      <span>{{
        Active !== (gamestate.Fighters.length !== 2)
          ? ''
          : (Active ? !White[0] : !!Vote)
            ? '✅'
            : '🕑'
      }}</span>
    </div>
  </div>
  <div v-if="playerIndex !== -1">
    <form
      v-if="!!gamestate.Streak || gamestate.Fighters.length === 2"
      id="fighters"
      @change="vote"
    >
      <label
        v-for="({ Black, White }, i) in gamestate.Fighters"
        :key="i"
        class="fighter"
        ><input
          v-model="optFighter"
          type="radio"
          name="fighter"
          :value="i + 1"
          :disabled="!!gamestate.Players[playerIndex].Active"
        />
        <span class="card white">{{ White }}</span>
        <span class="card black">{{ Black }}</span></label
      >
    </form>
    <form
      v-if="!!gamestate.Players[playerIndex].White[0]"
      id="hand"
      @submit="choose"
    >
      <fieldset>
        <label
          v-for="(text, i) in gamestate.Players[playerIndex].White"
          :key="i"
        >
          <input v-model="optWhite" type="radio" name="white" :value="i" />
          <span class="card white">{{ text }}</span></label
        >
      </fieldset>
      <fieldset>
        <label
          v-for="(text, i) in gamestate.Players[playerIndex].Black"
          :key="i"
        >
          <input v-model="optBlack" type="radio" name="black" :value="i" />
          <span class="card black">{{ text }}</span></label
        >
      </fieldset>
      <button type="submit" :disabled="!optWhite || !optBlack">Submit</button>
    </form>
  </div>
  <form v-if="playerIndex === -1" @submit="join">
    <label>
      Enter your name:
      <input v-model="playerName" type="text" />
    </label>
    <button type="submit">Join</button>
  </form>
</template>
