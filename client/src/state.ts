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

export interface GameSettings {
  Goal: number;
  HandSize: number;
}

export interface GameState {
  Done: boolean;
  Players: Player[];
  Fighters: Fighter[];
  Settings: GameSettings;
  Streak: number;
}

export const defaultGameState: GameState = {
  Done: false,
  Players: [],
  Fighters: [],
  Settings: { Goal: 3, HandSize: 6 },
  Streak: 0,
};

export function canVote(
  { Fighters: [a, b] }: GameState,
  playerIndex: number,
): boolean {
  return a && b && a.Player !== playerIndex && b.Player !== playerIndex;
}

export function getWinners({
  Players,
  Settings: { Goal },
}: GameState): string[] {
  return Players.filter(({ Points }) => Points >= Goal).map(({ Name }) => Name);
}
