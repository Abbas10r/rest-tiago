package main

import (
	"net/http"
	"socialApp/internal/store"
)

type CreatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJson(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserId:  1,
	}

	if err := app.store.Posts.Create(r.Context(), &post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
