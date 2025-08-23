package main

type GameSettings struct {
	Goal          byte
	FighterBlacks byte
	HandWhites    byte
	HandBlacks    byte
}

func DefaultGameSettings() GameSettings {
	return GameSettings{
		Goal:          6,
		FighterBlacks: 1,
		HandWhites:    3,
		HandBlacks:    3,
	}
}

type Fighter struct {
	Player   byte
	White    string
	Black    []string
	Tiebreak string
}

type Player struct {
	Name   string
	Points byte
	Vote   byte
	White  []string
	Black  []string
}

func (player *Player) DrawWhite(n byte, deck *Deck) {
	for range n {
		player.White = append(player.White, deck.Draw())
	}
}

func (player *Player) DrawBlack(n byte, deck *Deck) {
	for range n {
		player.Black = append(player.Black, deck.Draw())
	}
}

func (player *Player) Play(index byte, white byte, black []byte) Fighter {
	fighter := Fighter{
		Player: index,
		White:  player.White[white-1],
		Black:  make([]string, 0, len(black)),
	}
	for _, choice := range black {
		fighter.Black = append(fighter.Black, player.Black[choice-1])
	}
	player.White = player.White[:0]
	player.Black = player.Black[:0]
	return fighter
}

func (player *Player) Reset() {
	player.Vote = 0
	player.Points = 0
	player.Black = player.Black[:0]
	player.White = player.White[:0]
}

type GameState struct {
	Done     bool
	Players  []Player
	Fighters []Fighter
	Settings GameSettings
	Streak   byte
	nextUp   byte
	white    Deck
	black    Deck
}

func NewGameState(cards *Cards) *GameState {
	return &GameState{
		Players:  make([]Player, 0, 6),
		Fighters: make([]Fighter, 0, 2),
		Settings: DefaultGameSettings(),
		white:    NewDeck(cards.white),
		black:    NewDeck(cards.black),
	}
}

func (state *GameState) Choose(player byte, white byte, black []byte) {
	state.Fighters = append(state.Fighters, state.Players[player].Play(player, white, black))
}

func (state *GameState) draw(player *Player) {
	player.DrawWhite(state.Settings.HandWhites, &state.white)
	player.DrawBlack(state.Settings.HandBlacks, &state.black)
}

func (state *GameState) Reset() {
	for i := range state.Players {
		state.Players[i].Reset()
	}
	state.Done = false
	state.Fighters = state.Fighters[:0]
	state.Streak = 0
	state.white.Shuffle()
	state.black.Shuffle()
	for i := range min(len(state.Players), 2) {
		state.draw(&state.Players[i])
	}
	state.nextUp = 0
}

func (state *GameState) AddPlayer(name string) int {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return i
		}
	}
	playerIndex := len(state.Players)
	state.Players = append(state.Players, Player{
		Name:  name,
		White: make([]string, 0, state.Settings.HandBlacks),
		Black: make([]string, 0, state.Settings.HandBlacks),
	})
	if playerIndex < 2 {
		state.advanceNextUp()
	}
	return playerIndex
}

func (state *GameState) SetGoal(goal byte) {
	state.Settings.Goal = goal
	state.Done = false
	for i := range state.Players {
		if state.Players[i].Points >= goal {
			state.Done = true
			return
		}
	}
}

func (state *GameState) SetFighterBlacks(count byte) {
	state.Settings.FighterBlacks = count
	countI := int(count)
	for i := range state.Fighters {
		if len(state.Fighters[i].Black) > countI {
			state.Fighters[i].Black = state.Fighters[i].Black[:count]
		}
	}
}

func (state *GameState) SetHandBlacks(count byte) {
	state.Settings.HandBlacks = count
	for i := range state.Players {
		player := &state.Players[i]
		hand := byte(len(player.Black))
		if hand == 0 {
			continue
		}
		if hand < count {
			player.DrawBlack(count-hand, &state.black)
			continue
		}
		if hand > count {
			player.Black = player.Black[:count]
		}
	}
}

func (state *GameState) SetHandWhites(count byte) {
	state.Settings.HandWhites = count
	for i := range state.Players {
		player := &state.Players[i]
		hand := byte(len(player.White))
		if hand == 0 {
			continue
		}
		if hand < count {
			player.DrawWhite(count-hand, &state.white)
			continue
		}
		if hand > count {
			player.White = player.White[:count]
		}
	}
}

func (state *GameState) incrementNextUp() {
	if state.nextUp == byte(len(state.Players)-1) {
		state.nextUp = 0
	} else {
		state.nextUp += 1
	}
}

func (state *GameState) advanceNextUp() {
	state.incrementNextUp()
	if len(state.Fighters) == 1 && state.Fighters[0].Player == state.nextUp {
		state.incrementNextUp()
	}
	state.draw(&state.Players[state.nextUp])
}

func (state *GameState) Vote(player byte, vote byte) bool {
	state.Players[player].Vote = vote
	newVotes := 0
	oldVotes := 0
	for i := range byte(len(state.Players)) {
		if (len(state.Fighters) > 1 && state.Fighters[1].Player == i) ||
			(len(state.Fighters) > 0 && state.Fighters[0].Player == i) {
			continue
		}
		vote := state.Players[i].Vote
		if vote == 0 {
			return false
		}
		if vote == 1 {
			oldVotes += 1
		} else {
			newVotes += 1
		}
	}

	for i := range state.Players {
		state.Players[i].Vote = 0
	}

	if oldVotes > newVotes {
		state.Streak += 1
	} else if newVotes > oldVotes {
		state.Fighters[0] = state.Fighters[1]
		state.Streak = 1
	} else {
		state.Fighters[0].Tiebreak = state.white.Draw()
		state.Fighters[1].Tiebreak = state.white.Draw()
		return true
	}

	winner := &state.Players[state.Fighters[0].Player]
	winner.Points += 1
	if winner.Points >= state.Settings.Goal {
		state.Done = true
	}

	if state.Streak == 3 {
		state.Fighters = state.Fighters[:0]
		state.Streak = 0
		state.advanceNextUp()
	} else {
		state.Fighters = state.Fighters[:1]
		state.Fighters[0].Tiebreak = ""
	}

	state.advanceNextUp()
	return true
}
