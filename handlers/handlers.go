package handlers

import (
	"net/http"

	"github.com/GameStatisticAnalyst/ML-BE/middleware"
	"github.com/GameStatisticAnalyst/ML-BE/models"
	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello",
	})
}

func LoginHandler(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	
	// TODO : Validate user login
	user, err := models.GetUserByEmail(input.Email)

	ok, err := user.ComparePassword(input.Password)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "email or password salah",
		})
		return
	}

    accessToken, err := middleware.GenerateAccessToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
        return
    }

    refreshToken, err := middleware.GenerateRefreshToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}

func Refresh(c *gin.Context) {
    var input struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    claims, err := middleware.ValidateToken(input.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
        return
    }

    newAccessToken, err := middleware.GenerateAccessToken(claims.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": newAccessToken,
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
