package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
)

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	thread := models.Thread{}
	_ = thread.UnmarshalJSON(body)

	user, err := database.GetUserByNickname(thread.Author)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Author = user.Nickname

	forumSlug := mux.Vars(r)["slug"]
	forum, err := database.GetForumBySlug(forumSlug)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Forum = forum.Slug

	if thread.Slug != "" {
		existingThread, e := database.GetThreadBySlug(thread.Slug)
		if e == nil {
			network.WriteResponse(w, http.StatusConflict, existingThread)
			return
		}
	}

	e := database.CreateThread(&thread)
	if e != nil {
		network.WriteErrorResponse(w, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   e.Error(),
		})
		return
	}

	network.WriteResponse(w, http.StatusCreated, thread)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
