package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/service"
)

func main() {
	r := gin.Default()

	cityService := service.NewCityService()

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
	r.DELETE("/city/:uuid", func(c *gin.Context) {

		uuid := c.Param("uuid")

		err := cityService.Delete(uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.Run(":8081")
}
