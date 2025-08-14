package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type JsonCache struct {
	data []byte
	etag string
}

func NewJsonCache(v any) JsonCache {
	data, _ := json.Marshal(v)
	now, _ := time.Now().MarshalBinary()
	timestamp := sha256.Sum256(now)

	return JsonCache{
		data: append(data, byte('\n')),
		etag: fmt.Sprintf("\"%x\"", timestamp),
	}
}

func (j *JsonCache) Write(w http.ResponseWriter, r *http.Request) {
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
