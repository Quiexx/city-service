package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/repository"
	"github.com/Quiexx/city-service/internal/service"
)

func main() {

	// Load env

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

	// Connect to DB

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create kafka writer

	kafkaWriter := kafka.Writer{
		Addr:                   kafka.TCP(kafkaHost),
		Topic:                  "city-updates",
		AllowAutoTopicCreation: true,
	}

	cityRep := repository.NewPGCityRepository(db)
	cityService := service.NewFirstCityService(cityRep, &kafkaWriter)

	// REST API

	r := gin.Default()

	// Create new city
	r.POST("/city", func(c *gin.Context) {

		var city model.CreateCityRequest

		err := c.ShouldBindJSON(&city)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uuid, err := cityService.Create(city.Name, city.Population)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": uuid})
	})

	// Update city
	r.POST("/city/:id", func(c *gin.Context) {

		var city model.UpdateCityRequest

		id := c.Param("id")

		err := c.ShouldBindJSON(&city)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = cityService.Update(&model.City{ID: id, Name: city.Name, Population: city.Population})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	// Delete city
	r.DELETE("/city/:id", func(c *gin.Context) {

		uuid := c.Param("id")

		err := cityService.Delete(uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

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

	r.Run(":8081")
}
