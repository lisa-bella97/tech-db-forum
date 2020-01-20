package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	args := r.URL.Query()
	related := strings.Split(args.Get("related"), ",")

	result, err := database.GetPostFull(int64(id), related)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, result)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	postUpdate := models.PostUpdate{}
	_ = postUpdate.UnmarshalJSON(body)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	post, err := database.UpdatePost(&postUpdate, int64(id))
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	network.WriteResponse(w, http.StatusOK, post)
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	slugOrId := mux.Vars(r)["slug_or_id"]
	thread, err := database.GetThreadBySlug(slugOrId)
	if err != nil {
		id, _ := strconv.Atoi(slugOrId)
		thread, err = database.GetThreadById(int32(id))
		if err != nil {
			network.WriteErrorResponse(w, err)
			return
		}
	}

	posts := models.Posts{}
	_ = posts.UnmarshalJSON(body) // TODO: обрабатывать ошибки
	if len(posts) == 0 {
		network.WriteResponse(w, http.StatusCreated, posts)
		return
	}

	_, err = database.GetUserByNickname(posts[0].Author)
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	created := time.Now()

	for i := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum
		posts[i].Created = created
		err = database.CreatePost(&posts[i])
		if err != nil {
			network.WriteErrorResponse(w, err)
			return
		}
	}

	err = database.UpdateForumPosts(thread.Forum, len(posts))
	if err != nil {
		network.WriteErrorResponse(w, err)
		return
	}

	database.UpdateForumUsers(posts)

	network.WriteResponse(w, http.StatusCreated, posts)
}
