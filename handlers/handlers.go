package handlers

import (
	"net/http"

	"github.com/GameStatisticAnalyst/ML-BE/models"
	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello",
	})
}

func RegisterHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := models.NewUser(input.Username, input.Email, input.Password)

	isValid, errs := user.ValidateRegisterUser()
	if isValid == false {
		validationErrors := make(map[string]string)
		for field, err := range errs {
			validationErrors[field] = err.Error()
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "create user succeddddd",
	})
}
