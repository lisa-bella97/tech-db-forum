package handlers

import (
	"github.com/gorilla/mux"
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

	user, err := database.GetUserByNickname(forum.User)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	forum.User = user.Nickname

	existingForum, err := database.GetForumBySlug(forum.Slug)
	if err == nil {
		network.WriteResponse(w, http.StatusConflict, existingForum)
		return
	}

	e := database.CreateForum(forum)
	if e != nil {
		network.WriteErrorResponse(w, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   e.Error(),
		})
		return
	}

	network.WriteResponse(w, http.StatusCreated, forum)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	forum, err := database.GetForumBySlug(slug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	network.WriteResponse(w, http.StatusOK, forum)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
