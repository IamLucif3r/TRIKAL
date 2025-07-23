package main

import (
	"github.com/IamLucif3r/trikal/db"
	"github.com/IamLucif3r/trikal/feeds"
)

func main() {
	connStr := "postgres://username:password@localhost:5432/rssintel?sslmode=disable"
	db.InitDB(connStr)

	feeds.FetchAndStoreArticles()
}
