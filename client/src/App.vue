<script setup lang="ts">
import { computed, ref } from 'vue';
import { useApi } from './api';
import RocketIcon from './assets/rocket.svg?component';
import TrophyIcon from './assets/trophy.svg?component';
import GameSettingsForm from './components/GameSettingsForm.vue';
import JoinForm from './components/JoinForm.vue';
import PlayerInfo from './components/PlayerInfo.vue';
import { type GameSettings, canVote, getWinners } from './state';

const playerName = ref(localStorage.getItem('name') ?? '');
const playerIndex = ref(-1);
const optWhite = ref(0);
const optBlack = ref(0);
const optBlacks = ref<number[]>([]);
const optFighter = ref(0);

const { callApi, gamestate } = useApi(reset => {
  switch (reset) {
    case 'game':
      optWhite.value = 0;
      optBlack.value = 0;
      optBlacks.value.length = 0;
    case 'votes': //fallthrough
      optFighter.value = 0;
  }
});

const player = computed(() => gamestate.value.Players[playerIndex.value]);

const disableSubmit = computed(() => {
  if (!optWhite.value) {
    return true;
  }
  const fighterBlacks = gamestate.value.Settings.FighterBlacks;
  const black = optBlack.value;
  const blacks = optBlacks.value;
  if (fighterBlacks === 1) {
    return black === 0;
  }
  return blacks.length !== fighterBlacks;
});

function join(): void {
  const name = playerName.value;
  localStorage.setItem('name', name);

  callApi('PUT', '/api/join', { name })
    .then(data => (playerIndex.value = data.playerIndex))
    .catch(console.error);
}

function choose(): void {
  const isSingleChoice = gamestate.value.Settings.FighterBlacks === 1;

  callApi('PUT', '/api/choose', {
    player: playerIndex.value,
    white: optWhite.value,
    black: isSingleChoice ? [optBlack.value] : optBlacks.value,
  }).catch(console.error);
}

function patchGameSettings(event: Event): void {
  const setting = (event.target as HTMLInputElement).name as keyof GameSettings;

  callApi('PATCH', '/api/game', {
    [setting]: gamestate.value.Settings[setting],
  }).catch(console.error);
}

function vote(): void {
  callApi('PUT', '/api/vote', {
    player: playerIndex.value,
    fighter: optFighter.value,
  }).catch(console.error);
}

function resetGame(): void {
  callApi('PUT', '/api/reset').catch(console.error);
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
    <GameSettingsForm
      id="settings"
      v-model="gamestate.Settings"
      @change.prevent="patchGameSettings"
    />
  </div>
  <ul v-if="gamestate.Done" id="victory">
    <TrophyIcon />
    <li v-for="winner in getWinners(gamestate)" :key="winner">{{ winner }}</li>
    <button @click.prevent="resetGame">Start New Game</button>
  </ul>
  <div v-else-if="!!player" id="game">
    <form
      v-if="!!gamestate.Streak || gamestate.Fighters.length === 2"
      id="fighters"
      @change.prevent="vote"
    >
      <label
        v-for="({ Black, White, Tiebreak }, i) in gamestate.Fighters"
        :key="i"
        class="fighter"
      >
        <input
          v-model.number="optFighter"
          type="radio"
          name="fighter"
          :value="i + 1"
          :disabled="!canVote(gamestate, playerIndex)"
        />
        <span class="card white">{{ White }}</span>
        <span v-for="text in Black" :key="text" class="card black">
          {{ text }}
        </span>
        <span v-if="!!Tiebreak" class="card white">{{ Tiebreak }}</span>
      </label>
    </form>
    <form v-if="!!player.White.length" id="hand" @submit.prevent="choose">
      <button type="submit" :disabled="disableSubmit">
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
              v-if="gamestate.Settings.FighterBlacks > 1"
              v-model.number="optBlacks"
              type="checkbox"
              name="black"
              :value="i + 1"
            />
            <input
              v-else
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
  <JoinForm v-else id="join" v-model="playerName" @submit.prevent="join" />
</template>
