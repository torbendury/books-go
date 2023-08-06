// Package main only contains the main function, which starts up the server.
package main

import (
	"database/sql"
	"flag"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/torbendury/books-go/api"
	"github.com/torbendury/books-go/storage"
)

// main acts as a firestarter and configuration holder for the web server application.
// Besides creating a new server and starting it up, no actual logic should be placed here.
func main() {
	postgresMode := flag.Bool("postgres", false, "start server in postgres mode")

	postgresHost := flag.String("pghost", "localhost", "hostname (or IP) of postgres DB - only in postgres mode")
	postgresPort := flag.Int("pgport", 5432, "port of postgres DB - only in postgres mode")
	postgresUser := flag.String("pguser", "postgres", "user for postgres DB - only in postgres mode")
	postgresPass := flag.String("pgpass", "changeme", "password for postgres DB - only in postgres mode")
	postgresDb := flag.String("pgdatabase", "postgres", "database name of postgres DB - only in postgres mode")

	flag.Parse()

	var server *api.Server
	if *postgresMode {
		db := storage.OpenDB(*postgresHost, *postgresPort, *postgresUser, *postgresPass, *postgresDb)
		server = api.NewServer(storage.NewPostgresqlStorage(db), ":3000", fiber.Config{
			ServerHeader:          "books-go 0.0.1",
			AppName:               "books-go 0.0.1",
			IdleTimeout:           time.Duration(time.Second * 5),
			ReadTimeout:           time.Duration(time.Second * 5),
			DisableStartupMessage: true,
		})
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)
	} else {
		server = api.NewServer(storage.NewInMemoryStorage(), ":3000", fiber.Config{
			ServerHeader: "books-go 0.0.1-inmem-test",
			AppName:      "books-go 0.0.1-inmem-test",
			IdleTimeout:  time.Duration(time.Second),
			ReadTimeout:  time.Duration(time.Second),
		})
	}

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
