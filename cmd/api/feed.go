package main

import "net/http"

func (app *Application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination, filters

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(1))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
