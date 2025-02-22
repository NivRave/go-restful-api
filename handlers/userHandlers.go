package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"restful_api/cache"
	"restful_api/user"

	"github.com/asdine/storm/v3"
	"gopkg.in/mgo.v2/bson"
)

/*
** Delete an item **
Method - DELETE
Target - item
Endpoint - DELETE /collection/{id}
Request content - none
Successful response - 200 OK
Missing resource - 404 not found
*/
func UsersDeleteOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	err := user.DeleteOne(id)
	if err != nil {
		if err == storm.ErrNotFound {
			PostError(w, http.StatusNotFound)
			return
		}
		PostError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(r))
	w.WriteHeader(http.StatusOK)
}

/*
** Update an item **
Method - PATCH
Target - item
Endpoint - PATCH /collection/{id}
Request content - partial item data
Successful response - 200 OK + new item data
Missing resource - 404 not found
*/
func UsersPatchOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.GetOne(id)
	if err != nil {
		if err == storm.ErrNotFound {
			PostError(w, http.StatusNotFound)
			return
		}
		PostError(w, http.StatusInternalServerError)
		return
	}
	err = bodyToUser(r, u)
	if err != nil {
		PostError(w, http.StatusBadRequest)
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			PostError(w, http.StatusBadRequest)
		} else {
			PostError(w, http.StatusInternalServerError)
		}
		return
	}
	cache.Drop("/users")
	cw := cache.NewWriter(w, r)
	PostBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})
}

/*
** Replace an item **
Method - PUT
Target - item
Endpoint - PUT /collection/{id}
Request content - full item data
Successful response - 200 OK + new item data
Missing resource - 404 not found
*/
func UsersPutOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		PostError(w, http.StatusBadRequest)
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			PostError(w, http.StatusBadRequest)
		} else {
			PostError(w, http.StatusInternalServerError)
		}
		return
	}
	cache.Drop("/users")
	cw := cache.NewWriter(w, r)
	PostBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})
}

/*
** Create an item **
Method - POST
Target - collection
Endpoint - POST /collection
Request content - full item data
Successful response - 201 created + location
*/
func UsersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		PostError(w, http.StatusBadRequest)
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			PostError(w, http.StatusBadRequest)
		} else {
			PostError(w, http.StatusInternalServerError)
		}
		return
	}
	cache.Drop("/users")
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)
}

/*
** Access the collection **
Method - GET
Target - collection
Endpoint - GET /collection
Request content - none
Successful response - 200 OK + collection contents
*/
func UsersGetAll(w http.ResponseWriter, r *http.Request) {
	if cache.Serve(w, r) {
		return
	}
	u, err := user.GetAll()
	if err != nil {
		PostError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		PostBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	cw := cache.NewWriter(w, r)
	PostBodyResponse(cw, http.StatusOK, jsonResponse{"users": u})
}

/*
** Access an item **
Method - GET
Target - item
Endpoint - GET /collection/{id}
Request content - none
Successful response - 200 OK + item data
Missing resource - 404 not found
*/
func UsersGetOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	if cache.Serve(w, r) {
		return
	}
	u, err := user.GetOne(id)
	if err != nil {
		if err == storm.ErrNotFound {
			PostError(w, http.StatusNotFound)
			return
		}
		PostError(w, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodHead {
		PostBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	cw := cache.NewWriter(w, r)
	PostBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})
}

// Assisting method to parse a request's body to a User, if valid
func bodyToUser(r *http.Request, u *user.User) error {
	if r == nil {
		return errors.New("a request is required")
	}
	if r.Body == nil {
		return errors.New("request body is empty") // Extract to models / errors.go
	}
	if u == nil {
		return errors.New("a user is required")
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, u)
}
