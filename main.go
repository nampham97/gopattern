package main

import (
	"GoPattern/config"
	"GoPattern/db"
	"GoPattern/router"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("✅ Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", router.SetupRouter())
}
