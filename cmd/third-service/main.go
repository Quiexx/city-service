package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Quiexx/city-service/internal/service"
)

func main() {

	citiesUrl := "http://localhost:8082/city"
	cityService := service.NewThirdCityService(citiesUrl)

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

	r.Run(":8083")
}
