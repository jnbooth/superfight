export interface Fighter {
  Black: string;
  Player: number;
  White: string;
}

export interface Player {
  Active: boolean;
  Name: string;
  Points: number;
  Vote: number;
  White: string[];
  Black: string[];
}

export interface GameState {
  Players: Player[];
  Fighters: Fighter[];
  Streak: number;
  Tiebreak: string[];
}

export const defaultGameState: GameState = {
  Players: [],
  Fighters: [],
  Streak: 0,
  Tiebreak: [],
};
