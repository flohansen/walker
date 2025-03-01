package main

import (
	"flag"
	"log"

	"github.com/flohansen/walker/internal/application"
)

func main() {
	flags := application.Flags{}
	flag.IntVar(&flags.Server.ListenPort, "port", 3000, "The port to accept requests")
	flag.StringVar(&flags.Database.Host, "pg-host", "localhost", "The hostname of the PostgreSQL database")
	flag.IntVar(&flags.Database.Port, "pg-port", 5432, "The port of the PostgreSQL database")
	flag.StringVar(&flags.Database.Username, "pg-user", "", "The username for authenticate against the PostgreSQL database")
	flag.StringVar(&flags.Database.Password, "pg-pass", "", "The password for authenticate against the PostgreSQL database")
	flag.StringVar(&flags.Database.Database, "pg-db", "", "The database used to store data")
	flag.Parse()

	app := application.NewAPI(flags)
	if err := app.Run(application.SignalContext()); err != nil {
		log.Fatalf("api run error: %s", err)
	}
}
