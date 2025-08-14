package main

type Fighter struct {
	Player   byte
	White    string
	Black    string
	Tiebreak string
}

type Player struct {
	Name   string
	Points byte
	Vote   byte
	White  [3]string
	Black  [3]string
}

func NewPlayer(name string) Player {
	return Player{
		Name: name,
	}
}

func (player *Player) Draw(white *Deck, black *Deck) {
	player.White = [3]string{white.Draw(), white.Draw(), white.Draw()}
	player.Black = [3]string{black.Draw(), black.Draw(), black.Draw()}
}

func (player *Player) Play(index byte, white int, black int) Fighter {
	fighter := Fighter{
		Black:  player.Black[black-1],
		Player: index,
		White:  player.White[white-1],
	}
	player.Black = [3]string{}
	player.White = [3]string{}
	return fighter
}

func (player *Player) Reset() {
	player.Vote = 0
	player.Points = 0
}

type GameState struct {
	Players  []Player
	Fighters []Fighter
	Streak   byte
	nextUp   byte
	black    Deck
	white    Deck
}

func NewGameState(cards *Cards) *GameState {
	return &GameState{
		Players:  make([]Player, 0, 6),
		Fighters: make([]Fighter, 0, 2),
		Streak:   0,
		white:    NewDeck(cards.white),
		black:    NewDeck(cards.black),
		nextUp:   0,
	}
}

func (state *GameState) Choose(player byte, white int, black int) {
	state.Fighters = append(state.Fighters, state.Players[player].Play(player, white, black))
}

func (state *GameState) Reset() {
	for i := range state.Players {
		state.Players[i].Reset()
	}
	state.Fighters = state.Fighters[:0]
	state.Streak = 0
	state.black.Shuffle()
	state.white.Shuffle()
	state.Players[0].Draw(&state.white, &state.black)
	state.Players[1].Draw(&state.white, &state.black)
	state.nextUp = 0
}

func (state *GameState) AddPlayer(name string) int {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return i
		}
	}
	playerIndex := len(state.Players)
	state.Players = append(state.Players, NewPlayer(name))
	if playerIndex < 2 {
		state.advanceNextUp()
	}
	return playerIndex
}

func (state *GameState) resetVotes() {
	for i := range state.Players {
		state.Players[i].Vote = 0
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
	state.Players[state.nextUp].Draw(&state.white, &state.black)
}

func (state *GameState) Vote(player byte, vote byte) {
	state.Players[player].Vote = vote
	newVotes := 0
	oldVotes := 0
	for i := range byte(len(state.Players)) {
		if state.Fighters[0].Player == i || state.Fighters[1].Player == i {
			continue
		}
		vote := state.Players[i].Vote
		if vote == 0 {
			return
		}
		if vote == 1 {
			oldVotes += 1
		} else {
			newVotes += 1
		}
	}

	state.resetVotes()

	if oldVotes > newVotes {
		state.Streak += 1
	} else if newVotes > oldVotes {
		state.Fighters[0] = state.Fighters[1]
		state.Streak = 1
	} else {
		state.Fighters[0].Tiebreak = state.white.Draw()
		state.Fighters[1].Tiebreak = state.white.Draw()
		for i := range state.Players {
			state.Players[i].Vote = 0
		}
		return
	}

	state.resetVotes()

	state.Players[state.Fighters[0].Player].Points += 1

	if state.Streak == 3 {
		state.Fighters = state.Fighters[:0]
		state.Streak = 0
		state.advanceNextUp()
	} else {
		state.Fighters = state.Fighters[:1]
		state.Fighters[0].Tiebreak = ""
	}

	state.advanceNextUp()
}
