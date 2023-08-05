// Package main only contains the main function, which starts up the server.
package main

import (
	"database/sql"
	"flag"
	"fmt"
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
		fmt.Println("Starting server in PostgreSQL mode")
		fmt.Println("DB Host: ", *postgresHost)
		fmt.Println("DB Port: ", *postgresPort)
		fmt.Println("DB User: ", *postgresUser)
		fmt.Println("DB Pass: ", *postgresPass)
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
			fmt.Println("Closing DB.")
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
	server.Start()
}
