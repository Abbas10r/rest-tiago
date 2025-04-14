package main

import (
	"net/http"
	"socialApp/internal/store"
	"strconv"

	"github.com/gorilla/mux"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  1, //Change later
	}

	if err := app.store.Posts.Create(r.Context(), &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "")
		return
	}
	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		app.notFound(w, r, err)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	comments, err := app.store.Comments.GetByPostID(ctx, idInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
