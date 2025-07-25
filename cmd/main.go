package main

import (
	"database/sql"
	"log"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/database"
	"github.com/iamlucif3r/trikal/internal/types"
	"github.com/iamlucif3r/trikal/pkg"
)

var Config *types.Config
var Db *sql.DB

func init() {
	log.Println("Initializing TRIKAL ...")
	log.Println("Initializing configuration...")
	Config = &types.Config{}

	err := config.SetConfig(Config)
	if err != nil {
		log.Printf("Error setting configuration: %v\n", err)
		return
	}
	Db, err = database.ConnectDB(*Config)
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return
	}
	log.Println("Configuration initialized successfully.")
}
func main() {
	log.Println("TRIKAL is running with the following configuration : ")

	err := pkg.FetchNews()
	if err != nil {
		log.Printf("Error fetching news: %v\n", err)
		return
	}
}
