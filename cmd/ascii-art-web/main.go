package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ascii-art-web/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(":" + port)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		if err := srv.Shutdown(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	fmt.Printf("Server starting on http://localhost%s\n", srv.Addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
