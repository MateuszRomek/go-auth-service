package main

import "net/http"

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

	session, err := app.store.Sessions.CreateSession(r.Context(), user.Id)

	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = WriteJsonResponse(w, http.StatusOK, session); err != nil {
		WriteJsonError(w, http.StatusBadRequest, "Failed to create session")
	}
}
