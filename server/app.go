package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	cards := LoadCards()
	gamestate := NewGameState(&cards)
	gameJson := NewJsonCache(&gamestate)

	http.HandleFunc("/api/poll", func(w http.ResponseWriter, r *http.Request) {
		gameJson.Write(w, r)
	})

	http.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		gamestate.AddPlayer(r.FormValue("name"))
		gameJson = NewJsonCache(&gamestate)
		gameJson.Write(w, r)
	})

	http.HandleFunc("/api/choose", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		white, _ := strconv.ParseInt(r.FormValue("white"), 10, 0)
		black, _ := strconv.ParseInt(r.FormValue("black"), 10, 0)
		gamestate.Choose(byte(player), int(white), int(black))
		gameJson = NewJsonCache(&gamestate)
		gameJson.Write(w, r)
	})

	http.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		vote, _ := strconv.ParseUint(r.FormValue("fighter"), 10, 8)
		gamestate.Vote(byte(player), byte(vote))
		gameJson = NewJsonCache(&gamestate)
		gameJson.Write(w, r)
	})

	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		gamestate.Reset()
		gameJson = NewJsonCache(&gamestate)
		gameJson.Write(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("../client/dist")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
