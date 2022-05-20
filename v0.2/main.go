package main

import (
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/config"
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/handlers"
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/load"
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Test data.
var (
	TestAreaCode        = "E05011362"
	TestAreaName        = "Disbury East"
	TestAreaProfileName = "Resident Population for Disbury East, Census 2021"
)

// Flags
var (
	fLoadFiles []string
)

func main() {
	if err := run(); err != nil {
		log.Err("application error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cmd := &cobra.Command{}
	cmd.AddCommand(initCMD(), apiCMD())

	return cmd.Execute()
}

func initCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initalise the database, drops existing tables/data, recreates the schema and populates with a default area",
		Long: `The init command re-initalises the area_profiles database. Any existing tables are dropped, recreated and populated with a default area.
Use the init command to create the database for the first time or to tear down and recreate an existing database from scratch.

Using the -l flag you can specify 1 or more data files to load. If no file(s) are specified the key stats tables will be empty.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Get()
			if err != nil {
				return err
			}

			db, err := store.New(cfg.Username, cfg.Password, cfg.Database)
			if err != nil {
				return err
			}

			defer db.Close()

			if err := db.Init(TestAreaCode, TestAreaName, TestAreaProfileName); err != nil {
				return err
			}

			if len(fLoadFiles) == 0 {
				log.Info("init completed successfully")
				return nil
			}

			log.Info("loading test data into area_profiles database")
			for _, f := range fLoadFiles {
				fName := filepath.Join("load", f)

				if err := load.DataFromFile(fName, db); err != nil {
					return err
				}

				log.Info("successfully loaded test data: %s", fName)
			}

			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&fLoadFiles, "load", "l", []string{}, "A list of data import files to load (Optional). Format -l=file1 -l=file2 -l=fileN")
	return cmd
}

func apiCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Start the demo area profiles API.",
		Long: `Start the demo area profiles API. The API runs on port :8080 and exposes the following endpoints:
	GET: /profiles
	GET: /profiles/{area_code}
	GET: /profiles/{area_code}/stats
	GET: /profiles/{area_code}/stats/versions
	GET: /profiles/{area_code}/stats/versions/{version}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Get()
			if err != nil {
				return err
			}

			db, err := store.New(cfg.Username, cfg.Password, cfg.Database)
			if err != nil {
				return err
			}

			sigChan := make(chan os.Signal, 0)
			signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

			go func() {
				s := <-sigChan
				log.Warn("signal received initiating shutdown: %+v\n", s)
				db.Close()
				os.Exit(0)
			}()

			r := handlers.Initalise(db)

			log.Info("api ready to receive requests port :8080")
			if err := http.ListenAndServe(":8080", r); err != nil {
				return errors.Wrap(err, "errot shutting down http server")
			}
			return nil
		},
	}
}
