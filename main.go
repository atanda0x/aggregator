package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/atanda0x/aggregator/api"
	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/rss"
	"github.com/atanda0x/aggregator/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	queries := sqlc.New(conn)
	srv := api.NewServer(*queries)

	// Start scraping in a separate goroutine
	go rss.StartScraping(queries, 10, time.Minute)

	err = srv.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
