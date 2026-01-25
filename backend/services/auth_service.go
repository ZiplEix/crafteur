package services

import (
	"errors"
	"time"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte("super_secret_key_crafteur")

func Register(username, password string) error {
	existingUser, _ := database.GetUserByUsername(username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &core.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return database.CreateUser(user)
}

func Login(username, password string) (string, error) {
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(SecretKey)
}
