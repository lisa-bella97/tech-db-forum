package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"io/ioutil"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, errCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)
	marshalBody, err := json.Marshal(errorResponse{Message: errMsg})
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func WriteResponse(w http.ResponseWriter, code int, body interface{ MarshalJSON() ([]byte, error) }) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	marshalBody, err := body.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	f := models.Forum{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = f.UnmarshalJSON(body)

	_, err := database.GetUserByNickname(f.User)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+f.User)
		return
	}

	existingForum, err := database.GetForumBySlug(f.Slug)
	if err == nil {
		WriteResponse(w, http.StatusConflict, existingForum)
		return
	}

	err = database.CreateForum(f)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, f)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
