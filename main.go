/*
 * main.go
 *
 * Note to myself - finish up the route "/api/shorten" take your time
 */

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goccy/go-json" // https://github.com/goccy/go-json#benchmarks
)

func init() {
	// try load the load environment variables if there is an error log it
	err := LoadEnvVars()
	if err != nil {
		log.Fatalf("[Error] Loading Environment Variables: %s", err)
	}
	
	// try to make a connection to the redis database if it fails log the error
	err = ConnectToRedis()
	if err != nil {
		log.Fatalf("[Error] Connecting To Redis DataBase: %s", err)
	}
	
	// will be used to help generate random string
	rand.Seed(time.Now().UnixNano())
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
		
		result, err := rdb.Get(ctx, urlParam).Result()
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
		key := RandomString(10)
		_, err := rdb.Set(ctx, key, body.Url, time.Hour*168).Result()
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
		
		// Log success
		log.Printf("URL shortened. Key: %s, Original URL: %s", key, body.Url)
		
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
