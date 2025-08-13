export interface Fighter {
  Player: number;
  White: string;
  Black: string;
  Tiebreak: string;
}

export interface Player {
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
}

export const defaultGameState: GameState = {
  Players: [],
  Fighters: [],
  Streak: 0,
};

export function canVote(
  { Fighters: [a, b] }: GameState,
  playerIndex: number,
): boolean {
  return a && b && a.Player !== playerIndex && b.Player !== playerIndex;
}
