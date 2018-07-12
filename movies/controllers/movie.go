package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"

	"github.com/coderminer/microservice/dao"
	"github.com/coderminer/microservice/helper"
	"github.com/coderminer/microservice/models"
)

const (
	db         = "Movie"
	collection = "MovieModel"
)

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request")
		return
	}
	movie.Id = bson.NewObjectId().Hex()
	if err := dao.Insert(db, collection, movie); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, movie)
}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie
	if err := dao.FindAll(db, collection, nil, nil, &movies); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var movie models.Movie
	if err := dao.FindOne(db, collection, bson.M{"_id": id}, nil, &movie); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, movie)
}
