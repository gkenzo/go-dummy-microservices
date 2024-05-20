package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	defaultLogServiceUrl = "http://logger-service"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
		app.errorJSON(w, errors.New("error reading body"))
		return
	}

	switch rPayload.Action {
	case "auth":
		app.authenticate(w, rPayload.Auth)
	case "log":
		app.log(w, rPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	baseAuthSvURL, exists := os.LookupEnv("AUTH_SERVICE_URL")
	if !exists {
		baseAuthSvURL = "http://authentication-service/authenticate"
	}

	req, err := http.NewRequest("POST", baseAuthSvURL, bytes.NewBuffer(jsonData))
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

		return
	default:
		app.errorJSON(w, errors.New(fmt.Sprint(res.StatusCode)))
		log.Println(res.StatusCode, res.Header, res.Body)
		return
	}
}

func (app *Config) log(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")

	if err != nil {
		log.Println("error while indenting json data")
		app.errorJSON(w, errors.New("error while parsing json data"))
		return
	}

	logServiceUrl, ok := os.LookupEnv("LOG_SERVICE_URL")
	if !ok {
		logServiceUrl = defaultLogServiceUrl
	}

	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
