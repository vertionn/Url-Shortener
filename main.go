/*
 * main.go
 *
 * Note to myself - finish up the route "/api/shorten" take your time
 */

package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goccy/go-json" // https://github.com/goccy/go-json#benchmarks
	"go.uber.org/zap"
)

func init() {

	// start the logger
	NewLogger()

	// check if we have the correct ? set if not log it
	err := ValidateEnvVars()
	if err != nil {
		// return since we logged the missing environment variables inside the func
		return
	}

	// try to make a connection to the redis database if it fails log the error
	err = ConnectToRedis()
	if err != nil {
		GetLogger().Errorw("Failed to connect to the redis database",
			zap.String("Error", err.Error()),
		)
		return
	}

	// will be used to help generate random string
	rand.NewSource(time.Now().UnixNano())

	GetLogger().Infow("Initialized finished", nil)
}

// ShortenUrlBody this will be used to decode the incoming body into
type ShortenUrlBody struct {
	Url string `json:"Url"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, "code")
		if urlParam == "" {
			err := WriteJSON(w, http.StatusNotFound, jsonStruct{
				ErrorMessage: "missing url code",
				Success:      false,
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		result, err := rdb.Get(context.Background(), urlParam).Result()
		if err != nil {
			// Return the JSON response and exit the function
			err := WriteJSON(w, http.StatusNotFound, jsonStruct{
				ErrorMessage: "Resource not found",
				Success:      false,
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// Process the result if the key was found
		http.Redirect(w, r, result, http.StatusFound)
		return
	})

	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		// here we are just hoping they are passing a valid url to the endpoint
		// Note to myself we need a validator thinking github.com/go-playground/validator

		var body ShortenUrlBody
		decoding := json.NewDecoder(r.Body)
		decoding.DisallowUnknownFields()

		if err := decoding.Decode(&body); err != nil {
			err := WriteJSON(w, http.StatusBadRequest, jsonStruct{
				Success:      false,
				ErrorMessage: "Invalid request body",
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// now we need to store the data in the redis db
		key := NewRandomString(10)
		_, err := rdb.Set(context.Background(), key, body.Url, time.Hour*168).Result()
		if err != nil {
			// Handle Redis set error
			err := WriteJSON(w, http.StatusInternalServerError, jsonStruct{
				Success:      false,
				ErrorMessage: "Error storing data in Redis",
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// Respond with success
		err = WriteJSON(w, http.StatusOK, jsonStruct{
			Success: true,
			Message: key,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
