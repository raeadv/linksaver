package handler

import (
	"fmt"
	"linksaver/server/database"
	"linksaver/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// LoginCredential is the expected payload for POST /api/login.
type LoginCredential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginCredential is the expected payload for POST /api/login.
type RegisterCredential struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func HandleLogin(gc *gin.Context) {
	var credential LoginCredential
	gc.ShouldBindJSON(&credential)

	user, err := gorm.G[database.User](database.DB).Where("username = ?", credential.Username).First(gc)
	if err != nil {
		gc.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Failed to get User credential, please check your username/password [1]",
		})
		return
	}

	err = utils.ValidatePassword(user.Password, credential.Password)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to validate User credential, please check your username/password [2]",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Username, user.Email)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to validate User credential, please try again [3]",
		})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status": true,
		"token":  token,
	})
}

func HandleRegister(gc *gin.Context) {
	var newUser database.User
	gc.ShouldBindJSON(&newUser)

	hashedPwd, err := utils.HashString(newUser.Password)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to register, try again later [1]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	newUser.Password = string(hashedPwd)

	gorm.WithResult()
	err = gorm.G[database.User](database.DB).Create(gc, &newUser)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create user, please try again later [2]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	token, err := utils.GenerateToken(newUser.ID.String(), newUser.Username, newUser.Email)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to generate access token, please try again in login page [3]",
		})
		return
	}

	gc.JSON(http.StatusCreated, gin.H{
		"status": true,
		"token":  token,
	})
}

func HandleRefreshToken(gc *gin.Context) {
	strToken := utils.GetAuthHeader(gc)
	if strToken == "" {
		gc.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "invalid token",
		})
		return
	}

	token, err := utils.ValidateToken(strToken)

	if err != nil {
		gc.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "failed to validate token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		gc.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "Token is invalid, please log in again",
		})
		return
	}

	newToken, err := utils.GenerateToken(claims["id"].(string), claims["username"].(string), claims["email"].(string))
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to generate new token",
		})
	}

	gc.JSON(http.StatusCreated, gin.H{
		"status": true,
		"token":  newToken,
	})

}
