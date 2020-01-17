package handlers

import (
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
	"strconv"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	posts := models.Posts{}
	_ = posts.UnmarshalJSON(body) // TODO: обрабатывать ошибки
	if len(posts) == 0 {
		network.WriteResponse(w, http.StatusCreated, posts)
		log.Print("len = 0")
		return
	}

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

	_, err = database.GetUserByNickname(posts[0].Author)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	for i := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum
		err = database.CreatePost(&posts[i])
		log.Print(posts[i])
		if err != nil {
			network.WriteErrorResponse(w, err)
			return
		}
	}

	network.WriteResponse(w, http.StatusCreated, posts)
}
