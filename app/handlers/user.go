package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	user := models.User{}
	_ = user.UnmarshalJSON(body)

	nickname := mux.Vars(r)["nickname"]
	users, err := database.GetUsersByNicknameOrEmail(nickname, user.Email)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(users) > 0 {
		network.WriteResponse(w, http.StatusConflict, users)
		return
	}

	user.Nickname = nickname
	err = database.CreateUser(user)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	network.WriteResponse(w, http.StatusCreated, user)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	user, err := database.GetUserByNickname(nickname)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+nickname)
		return
	}

	network.WriteResponse(w, http.StatusOK, user)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
