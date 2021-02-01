package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fikifh/fullstack/api/auth"
	"github.com/fikifh/fullstack/api/models"
	"github.com/fikifh/fullstack/api/responses"
	"github.com/fikifh/fullstack/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreatedPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) getPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	//check if the auth token is valid and get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	//if a user attempt to update a post not belongin to him
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Unauthorized"))
		return
	}

	//read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//start processing the request data
	postUpdated := models.Post{}
	err = json.Unmarshal(body, &postUpdated)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//also acheck of the request user id is equal to the one gotten from token
	if uid != postUpdated.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//is a valid post id given to us ?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//is this user authenticated ?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//check if the post is exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	//is the authenticated user, the owner of this post ?
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
