package main

import (
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

func main() {
	cfg := mysql.Cfg("flux-scribe-bot:us-west1:flux-scribe-database", "user", "")
	// cfg.DBName = "flux-scribe-database"
	cfg.ParseTime = true

	const timeout = 10 * time.Second
	cfg.Timeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout

	db, err := mysql.DialCfg(cfg)
	if err != nil {
		panic("couldn't dial: " + err.Error())
	}
	// Close db after this method exits since we don't need it for the
	// connection pooling.
	defer db.Close()

	var now time.Time
	fmt.Println(db.QueryRow("SELECT NOW()").Scan(&now))
	fmt.Println(now)
}
