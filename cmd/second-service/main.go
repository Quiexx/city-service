package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/repository"
	"github.com/Quiexx/city-service/internal/service"
)

func main() {
	dbHost := os.Getenv("SERVICE_DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("SERVICE_DB_PORT")
	if len(dbPort) == 0 {
		dbPort = "5434"
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if len(dbUser) == 0 {
		dbPort = "postgres"
	}

	dbPass := os.Getenv("POSTGRES_PASSWORD")
	if len(dbPass) == 0 {
		dbPass = "postgres"
	}

	dbName := os.Getenv("POSTGRES_DB")
	if len(dbName) == 0 {
		dbName = "second_db"
	}

	kafkaHost := os.Getenv("SERVICE_KAFKA_HOST")
	if len(kafkaHost) == 0 {
		kafkaHost = "localhost:9092"
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaHost},
		Topic:          "city-updates",
		Partition:      0,
		CommitInterval: time.Second,
	})
	kafkaReader.SetOffsetAt(context.Background(), time.Now())

	cityRep := repository.NewPGCityRepository(db)
	cityService := service.NewSecondCityService(cityRep)

	go startListeningTopic(kafkaReader, cityService)

	r := gin.Default()

	// Get city
	r.GET("/city/:id", func(c *gin.Context) {

		uuid := c.Param("id")

		city, err := cityService.GetById(uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": city.ID, "name": city.Name, "population": city.Population})
	})

	r.Run(":8082")
}

func startListeningTopic(r *kafka.Reader, cityService service.ICityService) {
	defer r.Close()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var update model.CityUpdate
		err = json.Unmarshal(m.Value, &update)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received udpate: %v", update)

		if "INSERT" == update.Operation {
			cityService.CreateWithId(update.City.ID, update.City.Name, update.City.Population)
		}

		if "UPDATE" == update.Operation {
			cityService.Update(&update.City)
		}

		if "DELETE" == update.Operation {
			cityService.Delete(update.City.ID)
		}
	}

}
