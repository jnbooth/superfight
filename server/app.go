package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	cards := LoadCards()
	gamestate := NewGameState(&cards)
	mu := sync.Mutex{}
	hub := NewHub(&cards)
	go hub.Run()

	http.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		playerIndex := -1
		{
			mu.Lock()
			defer mu.Unlock()
			playerIndex = gamestate.AddPlayer(r.FormValue("name"))
		}
		hub.SendEvent("gameupdate", gamestate)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"playerIndex\": %v}\n", playerIndex)
	})

	http.HandleFunc("/api/choose", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		white, _ := strconv.ParseInt(r.FormValue("white"), 10, 0)
		black, _ := strconv.ParseInt(r.FormValue("black"), 10, 0)
		{
			mu.Lock()
			defer mu.Unlock()
			gamestate.Choose(byte(player), int(white), int(black))
		}
		hub.SendEvent("gameupdate", gamestate)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		vote, _ := strconv.ParseUint(r.FormValue("fighter"), 10, 8)
		votesReset := false
		{
			mu.Lock()
			defer mu.Unlock()
			votesReset = gamestate.Vote(byte(player), byte(vote))
		}
		if votesReset {
			hub.SendEvent("reset", "votes")
		}
		hub.SendEvent("gameupdate", gamestate)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		{
			mu.Lock()
			defer mu.Unlock()
			gamestate.Reset()
		}
		hub.SendEvent("reset", "game")
		hub.SendEvent("gameupdate", gamestate)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/api/game", func(w http.ResponseWriter, r *http.Request) {
		{
			mu.Lock()
			defer mu.Unlock()
			goal := r.FormValue("Goal")
			if goal != "" {
				goal, _ := strconv.ParseUint(goal, 10, 8)
				gamestate.SetGoal(byte(goal))
			}
			handSize := r.FormValue("HandSize")
			if handSize != "" {
				handSize, _ := strconv.ParseUint(handSize, 10, 8)
				gamestate.SetHandSize(byte(handSize))
			}
		}
		hub.SendEvent("gameupdate", gamestate)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		client := NewClient()
		hub.Register(client)
		defer hub.Unregister(client)
		client.Run(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("../client/dist")))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
