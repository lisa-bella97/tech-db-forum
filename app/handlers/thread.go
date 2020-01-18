package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
	"strconv"
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
		existingThread, err := database.GetThreadBySlug(thread.Slug)
		if err == nil {
			network.WriteResponse(w, http.StatusConflict, existingThread)
			return
		}
	}

	err = database.CreateThread(&thread)
	if err != nil {
		network.WriteErrorResponse(w, err)
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
	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThreadBySlug(slugOrId)
	if err != nil {
		id, _ := strconv.Atoi(slugOrId)
		thread, err = database.GetThreadById(id)
		if err != nil {
			network.WriteErrorResponse(w, err)
			return
		}
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	vote := models.Vote{}
	_ = vote.UnmarshalJSON(body)

	newVotes, err := database.Vote(vote, thread.Id)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}
	thread.Votes = newVotes

	network.WriteResponse(w, http.StatusOK, thread)
}
