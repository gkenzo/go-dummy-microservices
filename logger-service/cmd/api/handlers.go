package main

import (
	"net/http"

	"github.com/gkenzo/go-dummy-microservices/log-service/cmd/api/data"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var log jsonPayload
	_ = app.readJSON(w, r, log)

	entry := data.LogEntry{
		Name: log.Name,
		Data: log.Data,
	}

	err := app.Models.LogEntry.Insert(entry)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
