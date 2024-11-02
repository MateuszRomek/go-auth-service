package main

import (
	"net/http"

	"github.com/mateuszromek/auth/internal/auth"
)

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type LogoutPayload struct {
	Token string `json:"token" validate:"required,min=32"`
}

type TokenPayload struct {
	Token string `json:"token" validate:"required,min=32"`
}

type SessionTokenResponse struct {
	Session string `json:"session"`
	Token   string `json:"token"`
}

type OkResponse struct {
	Message string `json:"message"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	err := ReadJSON(w, r, &payload)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to parse payload")
		return
	}

	err = app.validator.Struct(payload)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.store.Users.CreateUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, token, err := app.store.Sessions.CreateSession(r.Context(), user.Id)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := SessionTokenResponse{
		Session: session.Id,
		Token:   token,
	}

	if err = WriteJsonResponse(w, http.StatusOK, response); err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to create session")
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	err := ReadJSON(w, r, &payload)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to parse payload")
		return
	}

	err = app.validator.Struct(payload)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.store.Users.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Cannot find user with provided email")
		return
	}

	isValid, err := auth.VerifyPassword(payload.Password, user.PasswordHash)

	if err != nil || !isValid {
		WriteJsonError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	err = app.store.Sessions.RemoveAllUserSessions(r.Context(), user.Id)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to create session")
		return
	}

	session, token, err := app.store.Sessions.CreateSession(r.Context(), user.Id)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to create session")
		return
	}

	response := SessionTokenResponse{
		Session: session.Id,
		Token:   token,
	}

	WriteJsonResponse(w, http.StatusOK, response)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	var payload LogoutPayload
	err := ReadJSON(w, r, &payload)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to parse payload")
		return
	}

	err = app.validator.Struct(payload)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := app.store.Sessions.GetSessionByToken(r.Context(), payload.Token)

	if err != nil || session == nil {
		WriteJsonError(w, http.StatusNotFound, "Failed to logout. User does not have active session")
		return
	}

	err = app.store.Sessions.RemoveAllUserSessions(r.Context(), session.UserId)

	if err != nil {
		WriteJsonError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	response := OkResponse{
		Message: "OK",
	}

	WriteJsonResponse(w, http.StatusOK, response)
}

func (app *application) validateToken(w http.ResponseWriter, r *http.Request) {
	var payload TokenPayload
	err := ReadJSON(w, r, &payload)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to parse payload")
		return
	}

	err = app.validator.Struct(payload)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, token, err := app.store.Sessions.ValidateSessionByToken(r.Context(), payload.Token)

	if err != nil {
		WriteJsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := SessionTokenResponse{
		Session: session.Id,
		Token:   token,
	}

	WriteJsonResponse(w, http.StatusOK, response)
}
