package routers

import (
	"net/http"
	"restful_api/handlers"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func UsersRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			handlers.UsersGetAll(w, r)
			return
		case http.MethodPost:
			handlers.UsersPostOne(w, r)
			return
		case http.MethodHead:
			handlers.UsersGetAll(w, r)
			return
		case http.MethodOptions:
			handlers.PostOptionsResponse(w, []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}, nil)
			return
		default:
			handlers.PostError(w, http.StatusMethodNotAllowed)
		}
	}

	path = strings.TrimPrefix(path, "/users/")
	if !bson.IsObjectIdHex(path) {
		handlers.PostError(w, http.StatusNotFound)
		return
	}

	id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:
		handlers.UsersGetOne(w, r, id)
		return
	case http.MethodPut:
		handlers.UsersPutOne(w, r, id)
		return
	case http.MethodPatch:
		handlers.UsersPatchOne(w, r, id)
		return
	case http.MethodDelete:
		handlers.UsersDeleteOne(w, r, id)
		return
	case http.MethodHead:
		handlers.UsersGetOne(w, r, id)
		return
	case http.MethodOptions:
		handlers.PostOptionsResponse(w, []string{http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions}, nil)
		return
	default:
		handlers.PostError(w, http.StatusMethodNotAllowed)
	}
}
