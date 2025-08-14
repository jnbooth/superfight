<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import RocketIcon from './assets/rocket.svg?component';
import TrophyIcon from './assets/trophy.svg?component';
import JoinForm from './components/JoinForm.vue';
import PlayerInfo from './components/PlayerInfo.vue';
import { type GameState, canVote, defaultGameState, getWinners } from './state';

const etag = ref('');
const gamestate = ref(defaultGameState);
const playerName = ref('');
const playerIndex = ref(-1);
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
  const state: GameState = await res.json();
  etag.value = res.headers.get('ETag') ?? etag.value;
  if (state.Streak === 0 && state.Fighters.length === 1) {
    const fighter = state.Fighters[0];
    fighter.Black = '';
    fighter.White = '';
  }
  gamestate.value = state;
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

async function setGoal(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PATCH', '/api/game', {
    Goal: gamestate.value.Settings.Goal.toString(),
  });
}

async function setHandSize(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PATCH', '/api/game', {
    HandSize: gamestate.value.Settings.HandSize.toString(),
  });
}

async function vote(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/vote', {
    player: playerIndex.value.toString(),
    fighter: gamestate.value.Players[playerIndex.value].Vote.toString(),
  });
}

async function resetGame(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/reset');
}

function useIntervalFn(cb: () => void, ms: number) {
  let int: ReturnType<typeof setInterval> | undefined;
  onMounted(() => (int = setInterval(cb, ms)));
  onUnmounted(() => int && clearInterval(int));
}

useIntervalFn(poll, 1000);

poll();
</script>

<template>
  <div id="sidebar">
    <span class="term">Players</span>
    <ol id="players">
      <PlayerInfo
        v-for="({ Name, Points, Vote, White }, i) in gamestate.Players"
        :key="Name"
        :waiting="!!White.length || (!Vote && canVote(gamestate, i))"
        :points="Points"
        :you="i === playerIndex"
      >
        {{ Name }}
      </PlayerInfo>
    </ol>
    <span class="term">Settings</span>
    <form id="settings">
      <label>
        <span>Goal</span>
        <input
          v-model.number="gamestate.Settings.Goal"
          type="number"
          name="Goal"
          min="1"
          max="255"
          @change="setGoal"
        />
      </label>
      <label>
        <span>Hand Size</span>
        <input
          v-model.number="gamestate.Settings.HandSize"
          type="number"
          name="HandSize"
          min="1"
          max="255"
          @change="setHandSize"
        />
      </label>
    </form>
  </div>
  <ul v-if="gamestate.Done" id="victory">
    <TrophyIcon />
    <li v-for="winner in getWinners(gamestate)" :key="winner">{{ winner }}</li>
    <button @click="resetGame">Start New Game</button>
  </ul>
  <div v-else-if="playerIndex !== -1" id="game">
    <form id="fighters" @change="vote">
      <label
        v-for="({ Black, White, Tiebreak }, i) in gamestate.Fighters"
        :key="i"
        class="fighter"
      >
        <input
          v-model.number="gamestate.Players[playerIndex].Vote"
          type="radio"
          name="fighter"
          :value="i + 1"
          :disabled="!canVote(gamestate, playerIndex)"
        />
        <span class="card white">{{ White }}</span>
        <span class="card black">{{ Black }}</span>
        <span v-if="!!Tiebreak" class="card white">{{ Tiebreak }}</span>
      </label>
    </form>
    <form
      v-if="!!gamestate.Players[playerIndex].White.length"
      id="hand"
      @submit="choose"
    >
      <button type="submit" :disabled="!optWhite || !optBlack">
        <RocketIcon /><span>Submit</span>
      </button>
      <div>
        <fieldset>
          <label
            v-for="(text, i) in gamestate.Players[playerIndex].White"
            :key="i"
          >
            <input
              v-model="optWhite"
              type="radio"
              name="white"
              :value="i + 1"
            />
            <span class="card white">{{ text }}</span>
          </label>
        </fieldset>
        <fieldset>
          <label
            v-for="(text, i) in gamestate.Players[playerIndex].Black"
            :key="i"
          >
            <input
              v-model="optBlack"
              type="radio"
              name="black"
              :value="i + 1"
            />
            <span class="card black">{{ text }}</span>
          </label>
        </fieldset>
      </div>
    </form>
  </div>
  <JoinForm v-else id="join" v-model="playerName" @submit="join" />
</template>
