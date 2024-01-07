/*
 * utils.go
 *
 * This file contains the functions to load the environment variables, write json back tot the user
 */

package main

import (
	"math/rand"
	"net/http"
	
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
)

// LoadEnvVars  loads environment variables from a .env file located outside the project folder.
func LoadEnvVars() error {
	
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(length int) string {
	b := make([]rune, length)
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

// WriteCustomJSON encodes and sends a JSON response based on the provided interface.
func WriteCustomJSON(w http.ResponseWriter, statusCode int, val interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		return err
	}
	
	return nil
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
