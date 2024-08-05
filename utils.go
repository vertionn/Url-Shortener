/*
 * utils.go
 *
 * This file contains the functions to load the environment variables, write json back tot the user
 */

package main

import (
	"errors"
	"math/rand"
	"net/http"
	"os"

	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

var (
	envs = []string{
		"REDIS_URI",
	}
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// LoadEnvVars  loads environment variables from a .env file located outside the project folder.
func ValidateEnvVars() error {

	var i int
	for _, key := range envs {
		q := os.Getenv(key)
		if q == "" {
			GetLogger().Errorw("missing environment variables",
				zap.String("Key", key),
			)
			i++
		}
	}

	if i > 1 {
		return errors.New("Missing environment variables")
	} else {
		return nil
	}
}

// generates a random string with lenght of l
func NewRandomString(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// the below code is for returning json

type jsonStruct struct {
	ErrorMessage string `json:"ErrorMessage,omitempty"`
	Message      string `json:"Message,omitempty"`
	Success      bool   `json:"Success"`
}

// WriteJSON encodes and sends a JSON response using the JSONResponse structure.
func WriteJSON(w http.ResponseWriter, statusCode int, response jsonStruct) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}
