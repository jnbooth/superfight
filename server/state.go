package main

type Fighter struct {
	Black  string
	Player byte
	White  string
}

type Player struct {
	Active bool
	Name   string
	Points byte
	Vote   byte
	White  [3]string
	Black  [3]string
}

func newPlayer(name string) Player {
	return Player{
		Name: name,
	}
}

func (player *Player) draw(white *Deck, black *Deck) {
	player.Active = true
	player.White = [3]string{white.draw(), white.draw(), white.draw()}
	player.Black = [3]string{black.draw(), black.draw(), black.draw()}
}

func (player *Player) play(index byte, white int, black int) Fighter {
	fighter := Fighter{
		Black:  player.Black[black-1],
		Player: index,
		White:  player.White[white-1],
	}
	player.Black = [3]string{}
	player.White = [3]string{}
	return fighter
}

func (player *Player) endTurn() {
	player.Active = false
	player.Vote = 0
}

func (player *Player) reset() {
	player.endTurn()
	player.Points = 0
}

type GameState struct {
	Players  []Player
	Fighters []Fighter
	Streak   byte
	Tiebreak [2]string
	nextUp   int
	black    Deck
	white    Deck
}

func newGameState(cards *Cards) *GameState {
	return &GameState{
		Players:  make([]Player, 0, 6),
		Fighters: make([]Fighter, 0, 2),
		Streak:   0,
		white:    newDeck(cards.white),
		black:    newDeck(cards.black),
		nextUp:   0,
	}
}

func (state *GameState) choose(player byte, white int, black int) {
	state.Fighters = append(state.Fighters, state.Players[player].play(player, white, black))
}

func (state *GameState) reset() {
	for i := range state.Players {
		state.Players[i].reset()
	}
	state.Fighters = state.Fighters[:0]
	state.Streak = 0
	state.Tiebreak[0] = ""
	state.Tiebreak[1] = ""
	state.black.shuffle()
	state.white.shuffle()
	state.Players[0].draw(&state.white, &state.black)
	state.Players[1].draw(&state.white, &state.black)
	state.nextUp = 0
}

func (state *GameState) addPlayer(name string) int {
	for i := range state.Players {
		if state.Players[i].Name == name {
			return i
		}
	}
	playerIndex := len(state.Players)
	state.Players = append(state.Players, newPlayer(name))
	if playerIndex < 2 {
		state.advanceNextUp()
	}
	return playerIndex
}

func (state *GameState) resetVotes() {
	for i := range state.Players {
		state.Players[i].Active = false
		state.Players[i].Vote = 0
	}
}

func (state *GameState) advanceNextUp() {
	if state.nextUp == len(state.Players)-1 {
		state.nextUp = 0
	} else {
		state.nextUp += 1
	}
	state.Players[state.nextUp].draw(&state.white, &state.black)
}

func (state *GameState) vote(player byte, vote byte) {
	state.Players[player].Vote = vote
	newVotes := 0
	oldVotes := 0
	for i := range state.Players {
		if state.Players[i].Active {
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
		state.Tiebreak = [2]string{state.white.draw(), state.white.draw()}
		for i := range state.Players {
			state.Players[i].Vote = 0
		}
		return
	}

	for i := range state.Players {
		state.Players[i].endTurn()
	}

	state.Players[state.Fighters[0].Player].Points += 1

	if state.Streak == 3 {
		state.Fighters = state.Fighters[:0]
		state.advanceNextUp()
	} else {
		state.Fighters = state.Fighters[:1]
	}

	state.advanceNextUp()
	state.Tiebreak = [2]string{}
}
