package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

var (
	store = make(map[string]string)
	mu    sync.RWMutex
)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Controlled CPU load: 10,000 iterations finishes in ~2-5ms per request
	hash := sha256.Sum256([]byte(req.URL))
	for i := 0; i < 10000; i++ {
		hash = sha256.Sum256(hash[:])
	}
	shortCode := hex.EncodeToString(hash[:])[:6]

	mu.Lock()
	store[shortCode] = req.URL
	mu.Unlock()

	resp := ShortenResponse{ShortURL: "/" + shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler)
	mux.HandleFunc("/healthz", healthHandler)

	// Define strict timeouts so connections do not hang open indefinitely
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	server.ListenAndServe()
}