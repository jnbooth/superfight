<script setup lang="ts">
import { computed, ref } from 'vue';
import { callApi, useServerGameState } from './api';
import RocketIcon from './assets/rocket.svg?component';
import TrophyIcon from './assets/trophy.svg?component';
import JoinForm from './components/JoinForm.vue';
import PlayerInfo from './components/PlayerInfo.vue';
import { canVote, getWinners } from './state';

const gamestate = useServerGameState();
const playerName = ref('');
const playerIndex = ref(-1);
const optWhite = ref(0);
const optBlack = ref(0);

const player = computed(() => gamestate.value.Players[playerIndex.value]);

async function join(event: Event): Promise<void> {
  event.preventDefault();
  const data = await callApi('PUT', '/api/join', {
    name: playerName.value,
  });
  playerIndex.value = data.playerIndex;
}

async function choose(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/choose', {
    player: playerIndex.value,
    white: optWhite.value,
    black: optBlack.value,
  });
}

async function setGoal(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PATCH', '/api/game', {
    Goal: gamestate.value.Settings.Goal,
  });
}

async function setHandSize(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PATCH', '/api/game', {
    HandSize: gamestate.value.Settings.HandSize,
  });
}

async function vote(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/vote', {
    player: playerIndex.value,
    fighter: gamestate.value.Players[playerIndex.value].Vote,
  });
}

async function resetGame(event: Event): Promise<void> {
  event.preventDefault();
  await callApi('PUT', '/api/reset');
}
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
  <div v-else-if="!!player" id="game">
    <form id="fighters" @change="vote">
      <label
        v-for="({ Black, White, Tiebreak }, i) in gamestate.Fighters"
        :key="i"
        class="fighter"
      >
        <input
          v-model.number="player.Vote"
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
    <form v-if="!!player.White.length" id="hand" @submit="choose">
      <button type="submit" :disabled="!optWhite || !optBlack">
        <RocketIcon /><span>Submit</span>
      </button>
      <div>
        <fieldset>
          <label v-for="(text, i) in player.White" :key="i">
            <input
              v-model.number="optWhite"
              type="radio"
              name="white"
              :value="i + 1"
            />
            <span class="card white">{{ text }}</span>
          </label>
        </fieldset>
        <fieldset>
          <label v-for="(text, i) in player.Black" :key="i">
            <input
              v-model.number="optBlack"
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
