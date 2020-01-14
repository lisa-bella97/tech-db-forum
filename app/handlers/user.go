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
		network.WriteErrorResponse(w, err)
		return
	}
	if len(users) > 0 {
		network.WriteResponse(w, http.StatusConflict, users)
		return
	}
	user.Nickname = nickname

	err = database.CreateUser(user)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusCreated, user)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	user, err := database.GetUserByNickname(nickname)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, user)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	user := models.User{}
	_ = user.UnmarshalJSON(body)

	if user.Email != "" {
		exUser, err := database.GetUserByEmail(user.Email)
		if err == nil {
			network.WriteErrorResponse(w, &models.ModelError{
				ErrorCode: http.StatusConflict,
				Message:   "This email is already registered by user: " + exUser.Nickname,
			})
			return
		}
	}

	nickname := mux.Vars(r)["nickname"]
	user.Nickname = nickname

	err := database.UpdateUser(&user)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, user)
}
