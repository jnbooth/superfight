<script setup lang="ts">
import { computed, ref } from 'vue';
import { useApi } from './api';
import RocketIcon from './assets/rocket.svg?component';
import TrophyIcon from './assets/trophy.svg?component';
import CardMultiSelect from './components/CardMultiSelect.vue';
import CardSelect from './components/CardSelect.vue';
import FighterSelect from './components/FighterSelect.vue';
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

const { callApi, gamestate, settings } = useApi(reset => {
  console.log('resetting', reset);
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
  const fighterBlacks = settings.value.FighterBlacks;
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
  const isSingleChoice = settings.value.FighterBlacks === 1;

  callApi('PUT', '/api/choose', {
    player: playerIndex.value,
    white: optWhite.value,
    black: isSingleChoice ? [optBlack.value] : optBlacks.value,
  }).catch(console.error);
}

function patchGameSettings(event: Event): void {
  const setting = (event.target as HTMLInputElement).name as keyof GameSettings;

  callApi('PATCH', '/api/game/settings', {
    [setting]: settings.value[setting],
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
      v-model="settings"
      @change.prevent="patchGameSettings"
    />
  </div>
  <ul v-if="gamestate.Done" id="victory">
    <TrophyIcon />
    <li v-for="winner in getWinners(gamestate, settings)" :key="winner">
      {{ winner }}
    </li>
    <button @click.prevent="resetGame">Start New Game</button>
  </ul>
  <div v-else-if="!!player" id="game">
    <FighterSelect
      v-if="!!gamestate.Streak || gamestate.Fighters.length === 2"
      id="fighters"
      v-model="optFighter"
      :fighters="gamestate.Fighters"
      :disabled="!canVote(gamestate, playerIndex)"
      @change.prevent="vote"
    />
    <form v-if="!!player.White.length" id="hand" @submit.prevent="choose">
      <button type="submit" :disabled="disableSubmit">
        <RocketIcon /><span>Submit</span>
      </button>
      <div>
        <CardSelect
          v-model.number="optWhite"
          :cards="player.White"
          name="white"
        />
        <CardMultiSelect
          v-if="settings.FighterBlacks > 1"
          v-model="optBlacks"
          :cards="player.Black"
          name="black"
        />
        <CardSelect
          v-else
          v-model="optBlack"
          :cards="player.Black"
          name="black"
        />
      </div>
    </form>
  </div>
  <JoinForm v-else id="join" v-model="playerName" @submit.prevent="join" />
</template>
