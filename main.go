/*
 * main.go
 *
 * Note to my self - finish up the route "/api/shorten" take your time
 */
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

}

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
		//
	})

	http.ListenAndServe(":8080", r)
}
