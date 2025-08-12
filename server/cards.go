package main

import (
	"math/rand"
	"os"
	"strings"
)

type Cards struct {
	black []string
	white []string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadCardsFromFile(filename string) []string {
	dat, err := os.ReadFile(filename)
	check(err)
	lines := string(dat[:len(dat)-1])
	return strings.Split(lines, "\n")
}

func loadCards() Cards {
	black := loadCardsFromFile("../cards/black.txt")
	white := loadCardsFromFile("../cards/white.txt")
	return Cards{black: black, white: white}
}

type Deck struct {
	i        int
	shuffled []int
	values   []string
}

func newDeck(values []string) Deck {
	deck := Deck{
		i:        0,
		shuffled: make([]int, len(values)),
		values:   values,
	}
	for i := range deck.shuffled {
		deck.shuffled[i] = i
	}
	deck.shuffle()
	return deck
}

func (deck *Deck) draw() string {
	card := deck.values[deck.shuffled[deck.i]]
	deck.i += 1
	if deck.i == len(deck.shuffled) {
		deck.shuffle()
	}
	return card
}

func (deck *Deck) shuffle() {
	deck.i = 0
	rand.Shuffle(len(deck.shuffled), func(i int, j int) {
		deck.shuffled[i], deck.shuffled[j] = deck.shuffled[j], deck.shuffled[i]
	})
}
