package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"shodan-proxy/internal/api"
	"shodan-proxy/internal/config"
)

func main() {
	// 加载主配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 加载 Shodan keys
	if err := config.LoadShodanKeys(); err != nil {
		log.Fatalf("Error loading Shodan keys: %v", err)
	}

	mux := api.SetupRoutes()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Server starting on :8080")
	log.Fatal(server.ListenAndServe())
}
