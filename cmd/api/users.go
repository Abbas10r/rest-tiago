package main

import (
	"net/http"
	"socialApp/internal/store"
	"strconv"

	"github.com/gorilla/mux"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
// @Summary Fethes a user profile
// @Description Fetches a user profile by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} store.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
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

// FollowUser godoc
//
// @Summary Follows a user
// @Description Follows a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {string} string "User followed"
// @Failure 400 {object} error "User not found"
// @Security ApiKeyAuth
// @Router /users/{id}/follow [put]
func (app *Application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)
	vars := mux.Vars(r)
	folowedId, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, folowedId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	vars := mux.Vars(r)
	unfolowedId, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, unfolowedId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromContext(r *http.Request) store.User {
	user, _ := r.Context().Value(userCtx).(store.User)
	return user
}
