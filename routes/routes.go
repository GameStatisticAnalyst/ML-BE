package routes

import (
	"github.com/GameStatisticAnalyst/ML-BE/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", handlers.IndexHandler)
	r.POST("/register", handlers.RegisterHandler)
}

