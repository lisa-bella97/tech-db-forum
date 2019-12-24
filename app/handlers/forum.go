package handlers

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
)

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	forum := models.Forum{}
	_ = forum.UnmarshalJSON(body)

	_, err := database.GetUserByNickname(forum.User)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+forum.User)
		return
	}

	existingForum, err := database.GetForumBySlug(forum.Slug)
	if err == nil {
		network.WriteResponse(w, http.StatusConflict, existingForum)
		return
	}

	err = database.CreateForum(forum)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	network.WriteResponse(w, http.StatusCreated, forum)
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
