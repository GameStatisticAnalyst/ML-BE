package routes

import (
	"net/http"

	"github.com/GameStatisticAnalyst/ML-BE/handlers"
	"github.com/GameStatisticAnalyst/ML-BE/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/refresh", handlers.Refresh)
	}


	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware())	
	{
		protected.GET("/whoaim", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"username": "test",
				"email": "testis@gmail.com",
			})
		})
	}

	r.GET("/", handlers.IndexHandler)
}

