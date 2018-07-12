package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/coderminer/microservice/dao"
	"github.com/coderminer/microservice/helper"
	"github.com/coderminer/microservice/models"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

const (
	db         = "User"
	collection = "UserModel"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := dao.Insert(db, collection, user); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, user)

}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var users []models.User
	if err := dao.FindAll(db, collection, nil, nil, &users); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, users)
}

func UserBooking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	name := params["name"]
	var user models.User
	if err := dao.FindOne(db, collection, bson.M{"_id": name}, nil, &user); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request")
		return
	}
	res, err := http.Get("http://127.0.0.1:8003/booking/" + name)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request by name "+name)
		return
	}

	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request of booking by name "+name)
		return
	}
	var booking models.Booking
	var resResult models.Result
	resResult.Name = name
	var resInfo models.ResultInfo

	if err := json.Unmarshal(result, &booking); err == nil {
		for _, book := range booking.Books {
			resInfo.Date = book.Date
			for _, movie := range book.Movies {
				res, err := http.Get("http://127.0.0.1:8001/movies/" + movie)
				if err == nil {
					result, err := ioutil.ReadAll(res.Body)
					if err == nil {
						var movie models.Movie
						if err := json.Unmarshal(result, &movie); err == nil {
							resInfo.Movies = append(resInfo.Movies, movie)
						}
					}
				}
			}
			resResult.Books = append(resResult.Books, resInfo)
		}
		helper.ResponseWithJson(w, http.StatusOK, resResult)
	} else {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request")
	}

}
