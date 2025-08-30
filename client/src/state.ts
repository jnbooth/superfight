export interface Fighter {
  Player: number;
  White: string;
  Black: string[];
  Tiebreak: string;
}

export interface Player {
  Name: string;
  Points: number;
  Vote: number;
  White: string[];
  Black: string[];
}

export interface GameSettings {
  Goal: number;
  FighterBlacks: number;
  HandWhites: number;
  HandBlacks: number;
}

export interface GameState {
  Done: boolean;
  Players: Player[];
  Fighters: Fighter[];
  Streak: number;
}

export const defaultGameState: GameState = {
  Done: false,
  Players: [],
  Fighters: [],
  Streak: 0,
};

export const defaultGameSettings: GameSettings = {
  Goal: 6,
  FighterBlacks: 1,
  HandWhites: 3,
  HandBlacks: 3,
};

export function canVote(
  { Fighters: [a, b] }: GameState,
  playerIndex: number,
): boolean {
  return a && b && a.Player !== playerIndex && b.Player !== playerIndex;
}

export function getWinners(
  { Players }: GameState,
  { Goal }: GameSettings,
): string[] {
  return Players.filter(({ Points }) => Points >= Goal).map(({ Name }) => Name);
}
