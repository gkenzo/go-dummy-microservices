package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "broker hit",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var rPayload RequestPayload

	err := app.readJSON(w, r, &rPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch rPayload.Action {
	case "auth":
		app.authenticate(w, rPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	authSvURL, exists := os.LookupEnv("AUTH_SERVICE_URL")
	if !exists {
		log.Println("couldn't get authentication service url from env.")
		authSvURL = "http://authentication-service/authenticate"
		log.Println("using default env instead: ", authSvURL)
	}
	req, err := http.NewRequest("POST", authSvURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusUnauthorized:
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	case http.StatusAccepted:
		var jsonAuthService jsonResponse
		err = json.NewDecoder(res.Body).Decode(&jsonAuthService)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		if jsonAuthService.Error {
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		app.writeJSON(w, http.StatusAccepted, jsonResponse{
			Error:   false,
			Message: "Authenticated!",
			Data:    jsonAuthService.Data,
		})

	default:
		app.errorJSON(w, errors.New("error"))
		log.Println(res.StatusCode, res.Header, res.Body)
		return
	}
}
