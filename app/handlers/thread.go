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

	_, err := database.GetUserByNickname(thread.Author)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+thread.Author)
		return
	}

	forumSlug := mux.Vars(r)["slug"]
	_, err = database.GetForumBySlug(forumSlug)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusNotFound, "Can't find forum with slug "+forumSlug)
		return
	}
	thread.Forum = forumSlug

	if thread.Slug != "" {
		existingThread, err := database.GetThreadBySlug(thread.Slug)
		if err == nil {
			network.WriteResponse(w, http.StatusConflict, existingThread)
			return
		}
	}

	err = database.CreateThread(&thread)
	if err != nil {
		network.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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
