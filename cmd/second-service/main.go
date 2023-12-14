package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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

	dsn := "host=localhost user=postgres password=postgres dbname=second_db port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "city-updates",
		Partition:      0,
		CommitInterval: time.Second,
	})

	kafkaWriter := kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "city-requests",
		AllowAutoTopicCreation: true,
	}

	cityRep := repository.NewPGCityRepository(db)
	cityService := service.NewSecondCityService(cityRep, &kafkaWriter)

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
