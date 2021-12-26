package main

import (
	"encoding/json"
	"net/http"

	"github.com/amtux/user-api-golang/pkg/models"
)

func (a *app) userSignup(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.logger.Errorw("Eror decoding signup request body", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	err = a.userDB.Create(&user)
	if err != nil {
		a.logger.Errorw("Error entering user into DB", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	jwtToken, err := a.jwt.Encode(user.Id)
	if err != nil {
		a.logger.Errorw("Error generating auth token", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	t := TokenResponse{Token: jwtToken}
	t.Send(w)
}

func (a *app) userLogin(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.logger.Errorw("Eror decoding signup request body", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	if user.Email == "" || user.Password == "" {
		f := Failure{http.StatusBadRequest, "Request missing email or password"}
		f.Send(w)
		return
	}
	userId, err := a.userDB.Auth(&user)
	// TODO check for  sql.ErrNoRows error and return "invalid credentials"
	if err != nil {
		f := Failure{http.StatusUnauthorized, err.Error()}
		f.Send(w)
		return
	}
	jwtToken, err := a.jwt.Encode(userId)
	if err != nil {
		a.logger.Errorw("Error generating auth token", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	t := TokenResponse{Token: jwtToken}
	t.Send(w)
}

func (a *app) putUser(w http.ResponseWriter, r *http.Request) {
	userId, err := a.jwt.Decode(r)
	if err != nil {
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
	reqUser := models.User{}
	err = json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		f := Failure{http.StatusBadRequest, "Error decoding signup request body"}
		f.Send(w)
		return
	}
	if reqUser.FirstName == "" && reqUser.LastName == "" {
		f := Failure{http.StatusBadRequest, "Need firstname or lastname to update for the user record"}
		f.Send(w)
		return
	}
	user, err := a.userDB.GetBasedOnId(userId)
	if err != nil {
		f := Failure{http.StatusInternalServerError, "Issues with fetching user"}
		f.Send(w)
		return
	}
	user.FirstName = reqUser.FirstName
	user.LastName = reqUser.LastName

	err = a.userDB.Update(user)
	if err != nil {
		a.logger.Errorw("Failed to update user with new details", err)
		f := Failure{http.StatusBadRequest, err.Error()}
		f.Send(w)
		return
	}
}
