package main

import (
	"calorie_tracker_backend/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"

	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.POST("/entry/create", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry/:id/", routes.GetEntryById)
	router.GET("/ingredient/:ingredient", routes.GetEntriesByIngredient)

	router.PUT("/entry/update/:id", routes.UpdateEntry)
	router.PUT("/ingredient/update/:id", routes.UpdateIngredient)

	router.DELETE("/entry/delete/:id", routes.DeleteEntry)
	router.GET("/", routes.Land)

	router.Run(":" + port)

}
