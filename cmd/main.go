package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to TRIKAL - Your Cybersecurity News Aggregator",
		})
	})
	router.GET("/fetch-news", func(c *gin.Context) {
		err := pkg.FetchNews()
		if err != nil {
			log.Printf("Error fetching news: %v\n", err)
			c.JSON(500, gin.H{"error": "Failed to fetch news"})
			return
		}
		c.JSON(200, gin.H{"message": "News fetched successfully"})
	})
	router.Run(":3333")

}
