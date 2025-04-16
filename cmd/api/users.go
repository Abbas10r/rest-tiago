package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "")
		return
	}
	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, id)
	if err != nil {
		app.notFound(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
