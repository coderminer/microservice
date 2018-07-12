package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coderminer/microservice/dao"
	"github.com/coderminer/microservice/helper"
	"github.com/coderminer/microservice/models"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

const (
	db         = "Booking"
	collection = "BookModel"
)

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, "invalid request")
		return
	}
	booking.Id = bson.NewObjectId().Hex()
	if err := dao.Insert(db, collection, booking); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, booking)
}

func GetAllBooking(w http.ResponseWriter, r *http.Request) {
	var bookings []models.Booking
	if err := dao.FindAll(db, collection, nil, nil, &bookings); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, bookings)
}

func GetBookByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	var booking models.Booking
	if err := dao.FindOne(db, collection, bson.M{"name": name}, nil, &booking); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, booking)
}
