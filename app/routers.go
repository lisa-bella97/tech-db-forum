/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package app

import (
	"github.com/gorilla/mux"
	"github.com/lisa-bella97/tech-db-forum/app/handlers"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		//handler = log.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		handlers.Clear,
	},

	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/create",
		handlers.ForumCreate,
	},

	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/details",
		handlers.ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/threads",
		handlers.ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/users",
		handlers.ForumGetUsers,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/{id}/details",
		handlers.PostGetOne,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/{id}/details",
		handlers.PostUpdate,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/create",
		handlers.PostsCreate,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		handlers.Status,
	},

	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/{slug}/create",
		handlers.ThreadCreate,
	},

	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/details",
		handlers.ThreadGetOne,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/posts",
		handlers.ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/details",
		handlers.ThreadUpdate,
	},

	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/vote",
		handlers.ThreadVote,
	},

	Route{
		"UserCreate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/create",
		handlers.UserCreate,
	},

	Route{
		"UserGetOne",
		strings.ToUpper("Get"),
		"/api/user/{nickname}/profile",
		handlers.UserGetOne,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/profile",
		handlers.UserUpdate,
	},
}
