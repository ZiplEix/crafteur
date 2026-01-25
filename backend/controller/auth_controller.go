package controller

import (
	"net/http"
	"time"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/ZiplEix/crafteur/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	req := new(core.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	token, err := services.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	// Secure should be true in production, but leaving it off for dev/local unless user specified.
	// user didn't specify Secure default, but said "Cookie HTTP-Only".

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Expire immediately
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

func Me(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving user"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}
