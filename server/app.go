package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CachedJson struct {
	data []byte
	etag string
}

func cacheJson(v any) CachedJson {
	data, _ := json.Marshal(v)
	now, _ := time.Now().MarshalBinary()
	timestamp := sha256.Sum256(now)

	return CachedJson{
		data: append(data, byte('\n')),
		etag: fmt.Sprintf("\"%x\"", timestamp),
	}
}

func (j *CachedJson) write(w http.ResponseWriter, r *http.Request) {
	ifNoneMatch := r.Header["If-None-Match"]
	if len(ifNoneMatch) > 0 && ifNoneMatch[0] == j.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", j.etag)
	w.WriteHeader(http.StatusOK)
	w.Write(j.data)
}

func main() {
	r := mux.NewRouter()

	cards := loadCards()
	gamestate := newGameState(&cards)
	gameJson := cacheJson(&gamestate)

	r.HandleFunc("/api/poll", func(w http.ResponseWriter, r *http.Request) {
		gameJson.write(w, r)
	}).Methods("GET")

	r.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		gamestate.addPlayer(r.FormValue("name"))
		gameJson = cacheJson(&gamestate)
		gameJson.write(w, r)
	}).Methods("PUT")

	r.HandleFunc("/api/choose", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		white, _ := strconv.ParseInt(r.FormValue("white"), 10, 0)
		black, _ := strconv.ParseInt(r.FormValue("black"), 10, 0)
		gamestate.choose(byte(player), int(white), int(black))
		gameJson = cacheJson(&gamestate)
		gameJson.write(w, r)
	}).Methods("PUT")

	r.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
		player, _ := strconv.ParseUint(r.FormValue("player"), 10, 8)
		vote, _ := strconv.ParseUint(r.FormValue("fighter"), 10, 8)
		gamestate.vote(byte(player), byte(vote))
		gameJson = cacheJson(&gamestate)
		gameJson.write(w, r)
	}).Methods("PUT")

	r.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		gamestate.reset()
		gameJson = cacheJson(&gamestate)
		gameJson.write(w, r)
	}).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/dist")))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
