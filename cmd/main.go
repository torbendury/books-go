// Package main only contains the main function, which starts up the server.
package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/torbendury/books-go/api"
	"github.com/torbendury/books-go/storage"
)

// main acts as a firestarter and configuration holder for the web server application.
// Besides creating a new server and starting it up, no actual logic should be placed here.
func main() {
	server := api.NewServer(storage.NewInMemoryStorage(), ":3000", fiber.Config{
		ServerHeader: "books-go 0.0.1",
		AppName:      "books-go 0.0.1",
		IdleTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		// DisableStartupMessage: true,
	})

	server.Start()
}
