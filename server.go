package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"
)

func SetupBackend(dir string) *Hub {
	cards := LoadCards(dir)
	gamestate := NewGameState(&cards)
	mu := sync.Mutex{}
	hub := NewHub(&cards)

	http.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name, err := FormString(r, "name")
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		playerIndex := -1
		{
			mu.Lock()
			defer mu.Unlock()
			playerIndex = gamestate.AddPlayer(name)
		}
		hub.SendEvent("gameupdate", gamestate)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"playerIndex\": %v}\n", playerIndex)
	})

	http.HandleFunc("/api/choose", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		player, err := FormByte(r, "player", 0, byte(len(gamestate.Players)-1))
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		white, err := FormByte(r, "white", 1, gamestate.Settings.HandWhites)
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		black, err := FormBytes(r, "black", 1, gamestate.Settings.HandBlacks, int(gamestate.Settings.FighterBlacks))
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		{
			mu.Lock()
			defer mu.Unlock()
			gamestate.Choose(player, white, black)
		}
		hub.SendEvent("gameupdate", gamestate)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		player, err := FormByte(r, "player", 0, byte(len(gamestate.Players)-1))
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		vote, err := FormByte(r, "fighter", 1, 2)
		if err != nil {
			WriteInvalid(w, err)
			return
		}
		votesReset := false
		{
			mu.Lock()
			defer mu.Unlock()
			votesReset = gamestate.Vote(player, vote)
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
		if r.Method == http.MethodGet {
			data, _ := json.Marshal(gamestate)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}
		r.ParseForm()
		{
			mu.Lock()
			defer mu.Unlock()
			if len(r.Form["Goal"]) > 0 {
				goal, err := FormByte(r, "Goal", 1, 255)
				if err != nil {
					WriteInvalid(w, err)
					return
				}
				gamestate.SetGoal(goal)
			}
			if len(r.Form["FighterBlacks"]) > 0 {
				fighterBlacks, err := FormByte(r, "FighterBlacks", 1, 255)
				if err != nil {
					WriteInvalid(w, err)
					return
				}
				gamestate.SetFighterBlacks(fighterBlacks)
			}
			if len(r.Form["HandBlacks"]) > 0 {
				handBlacks, err := FormByte(r, "HandBlacks", 1, 255)
				if err != nil {
					WriteInvalid(w, err)
					return
				}
				gamestate.SetHandBlacks(handBlacks)
			}
			if len(r.Form["HandWhites"]) > 0 {
				handWhites, err := FormByte(r, "HandWhites", 1, 255)
				if err != nil {
					WriteInvalid(w, err)
					return
				}
				gamestate.SetHandWhites(handWhites)
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

	return hub
}

func SetupFrontend(dir string) {
	http.Handle("/", http.FileServer(http.Dir(path.Join(dir, "static"))))
}
