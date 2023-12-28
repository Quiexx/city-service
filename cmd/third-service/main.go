package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Quiexx/city-service/internal/service"
)

func main() {

	// Load env

	citiesUrl := os.Getenv("CITY_URL")
	if len(citiesUrl) == 0 {
		citiesUrl = "http://localhost:8082/city"
	}
	cityService := service.NewThirdCityService(citiesUrl)

	// API

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

	// Get all cities
	r.GET("/city", func(c *gin.Context) {

		cities, err := cityService.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cities": cities})
	})

	r.Run(":8083")
}
