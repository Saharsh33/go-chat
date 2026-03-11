package main

import (
	"log"
	"net/http"
	"os"

	"chat-server/internal/config"
	chathttp "chat-server/internal/http"
	postgres "chat-server/internal/storage/db"
	"chat-server/internal/websocket"

	_ "github.com/lib/pq"
)

func main() {
	//logging format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 1. Load config
	cfg := config.Load()

	// 2. Connect to Postgres (ONLY here)
	db := postgres.NewDB(cfg.PostgresDSN)
	if db == nil {
		panic("DB connection failed!!")
	}
	log.Println("Connected to Postgres")

	// 3. Create message store
	store := postgres.NewStore(db)

	// 4. Create hub
	hub := websocket.NewHub(store)
	go hub.Run()

	// 5. Auth endpoints
	http.HandleFunc("/register", chathttp.NewRegisterHandler(store, cfg.JWTSecret))
	http.HandleFunc("/login", chathttp.NewLoginHandler(store, cfg.JWTSecret))

	// 6. WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r, cfg.JWTSecret)
	})

	// 7. Health check (optional but good)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 8. Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

