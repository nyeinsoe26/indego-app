package main

import (
	"flag"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/config"
	"github.com/nyeinsoe26/indego-app/internal/app/api"
	"github.com/nyeinsoe26/indego-app/internal/app/db"
	"github.com/nyeinsoe26/indego-app/internal/app/services"
	"github.com/nyeinsoe26/indego-app/internal/gateways"
)

// Initialize and start the cronjob to fetch data every hour using the handler
func startCronJob(handler *api.Handler) {
	go func() {
		for {
			// Call the core logic to fetch and store data
			err := handler.FetchAndStoreIndegoWeatherData()
			if err != nil {
				// Log the error in case of failure
				log.Printf("Cronjob failed: %v", err)
			} else {
				log.Println("Cronjob successfully fetched and stored data")
			}

			// Sleep for an hour before running again
			time.Sleep(1 * time.Hour)
		}
	}()
}

func main() {
	// Define the --config flag to get the config file path
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	// Load the configuration file
	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL connection
	database, err := db.NewPostgresDB(config.GetDatabaseConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	// Initialize Indego and Weather clients
	indegoClient := gateways.NewIndegoClient()
	weatherClient := gateways.NewWeatherClient()

	// Initialize services
	indegoService := services.NewIndegoService(indegoClient, database)
	weatherService := services.NewWeatherService(weatherClient, database)

	// Initialize handlers
	handler := api.NewHandler(indegoService, weatherService)

	// Start the cronjob to fetch and store data every hour
	startCronJob(handler)

	// Initialize Gin router
	router := gin.Default()

	// Register routes and pass the handler
	api.RegisterRoutes(router, handler)

	// Start the server
	err = router.Run(":" + config.AppConfig.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
