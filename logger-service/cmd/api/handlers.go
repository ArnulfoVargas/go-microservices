package main

import (
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload JSONPayload

	_ = app.ReadJson(w, r, &payload)

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error: false,
		Message: "Logged",
	}

	app.WriteJSON(w, http.StatusAccepted, resp)
}
