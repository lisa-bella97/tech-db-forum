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
		network.WriteErrorResponse(w, err)
		return
	}

	forumSlug := mux.Vars(r)["slug"]
	_, e := database.GetForumBySlug(forumSlug)
	if e != nil {
		network.WriteErrorResponse(w, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find forum with slug " + forumSlug,
		})
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

	e = database.CreateThread(&thread)
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
