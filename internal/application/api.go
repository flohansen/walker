package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/flohansen/walker/generated/database"
	"github.com/flohansen/walker/internal/controller"
	"github.com/flohansen/walker/internal/repository"
	localsql "github.com/flohansen/walker/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/lib/pq"
)

type ServerFlags struct {
	ListenPort int
}

type DatabaseFlags struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (f DatabaseFlags) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		f.Username, f.Password, f.Host, f.Port, f.Database)
}

type Flags struct {
	Server   ServerFlags
	Database DatabaseFlags
}

type API struct {
	server *http.Server
}

func NewAPI(flags Flags) *API {
	dsn := flags.Database.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("could not open postgres database: %s", err)
	}

	d, err := iofs.New(localsql.MigrationsFS, "migrations")
	if err != nil {
		log.Fatalf("could not create migration source: %s", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		log.Fatalf("could not create migration instance: %s", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("could not migrate: %s", err)
	}

	repo := &repository.RouteSQLRepository{
		Queries: database.New(db),
	}

	return &API{
		server: &http.Server{
			Handler: controller.NewRoutes(repo),
			Addr:    fmt.Sprintf(":%d", flags.Server.ListenPort),
		},
	}
}

func (a *API) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		a.server.Shutdown(context.Background())
	}()

	return a.server.ListenAndServe()
}

func SignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()

		<-sig
		os.Exit(1)
	}()

	return ctx
}
