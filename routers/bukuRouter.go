package routers

import (
	"project-buku/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	// Menggunakan middleware basic auth
	// Detail Feature Kategori
	router.GET("/api/categories", controllers.Auth(), controllers.GetCategories)
	router.POST("api/categories", controllers.Auth(), controllers.CreateCategories)
	router.GET("/api/categories/:id", controllers.Auth(), controllers.GetEachCategories)
	router.DELETE("/api/categories/:id", controllers.Auth(), controllers.DeleteCategories)
	router.GET("/api/categories/:id/books", controllers.Auth(), controllers.GetEachBookCategories)

	router.GET("/api/books", controllers.Auth(), controllers.GetBuku)
	router.POST("/api/books", controllers.Auth(), controllers.CreateBuku)
	router.GET("/api/books/:id", controllers.Auth(), controllers.GetEachBuku)
	router.DELETE("/api/books/:id", controllers.Auth(), controllers.DeleteBuku)

	return router
}
