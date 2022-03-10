package main

import (
	"flag"
	"github.com/ONSdigital/dp-area-profiles-design-spike/handlers"
	"github.com/ONSdigital/dp-area-profiles-design-spike/store"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Test data.
var (
	TestAreaCode        = "E05011362"
	TestAreaName        = "Disbury East"
	TestAreaProfileName = "Resident Population for Disbury East, Census 2021"
)

type config struct {
	drop     bool
	username string
	password string
	database string
}

func main() {
	if err := run(); err != nil {
		log.Err("application error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log.Init("area-profiles-spike")
	log.Info("starting applicaiton")

	cfg, err := getDBConfig()
	if err != nil {
		return err
	}

	areaStore, err := store.New(cfg.username, cfg.password, cfg.database)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 0)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-sigChan
		log.Warn("signal received initiating shutdown: %+v\n", s)

		areaStore.Close()
		os.Exit(0)
	}()

	if cfg.drop {
		if err = areaStore.Init(TestAreaCode, TestAreaName, TestAreaProfileName); err != nil {
			return err
		}
	}

	r, err := handlers.Initialise(areaStore)
	if err != nil {
		return err
	}

	log.Info("api ready to receive requests port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		return errors.Wrap(err, "errot shutting down http server")
	}

	return nil
}

func getDBConfig() (*config, error) {
	var (
		drop       bool
		dbUsername string
		dbPassword string
		dbName     string
	)

	flag.BoolVar(&drop, "drop", false, "if true drop and recreate the database schema")
	flag.StringVar(&dbUsername, "u", "", "db username ")
	flag.StringVar(&dbPassword, "p", "", "db password")
	flag.StringVar(&dbName, "db", "", "db name")
	flag.Parse()

	if dbUsername == "" {
		return nil, errors.New("expected flag -u database username")
	}

	if dbPassword == "" {
		return nil, errors.New("expected flag -p database password")
	}

	if dbName == "" {
		return nil, errors.New("expected flag -db database name")
	}

	return &config{
		drop:     drop,
		username: dbUsername,
		password: dbPassword,
		database: dbName,
	}, nil
}
