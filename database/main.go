package main

import (
	json "encoding/json"
	"fmt"
	io "io/ioutil"

	"github.com/coderminer/microservice/dao"
	"github.com/coderminer/microservice/models"
)

func main() {
	initUser("./user.json")
	initMovie("./movie.json")
	initShowtime("./showtimes.json")
	initBooking("./booking.json")
}

func initUser(filename string) {
	d, err := io.ReadFile(filename)
	if err != nil {
		panic("read user user.json error")
	}
	var users []models.User
	var ui []interface{}
	err = json.Unmarshal(d, &users)
	if err != nil {
		panic("unmarshal json file error")
	}
	for _, d := range users {
		ui = append(ui, d)
	}
	fmt.Println("users: ", ui)
	err = dao.Insert("User", "UserModel", ui...)
	if err != nil {
		panic(fmt.Sprintf("insert user %s", err))
	}

}

func initMovie(filename string) {
	d, _ := io.ReadFile(filename)
	var movies []models.Movie
	json.Unmarshal(d, &movies)
	var inter []interface{}

	for _, item := range movies {
		inter = append(inter, item)
	}
	err := dao.Insert("Movie", "MovieModel", inter...)
	if err != nil {
		fmt.Println("insert movie error,", err)
	}
}

func initShowtime(filename string) {
	d, _ := io.ReadFile(filename)
	var shows []models.ShowTimes
	json.Unmarshal(d, &shows)
	var inter []interface{}
	for _, item := range shows {
		inter = append(inter, item)
	}
	err := dao.Insert("ShowTimes", "ShowModel", inter...)
	if err != nil {
		fmt.Println("insert ShowTimes error,", err)
	}
}

func initBooking(filename string) {
	d, _ := io.ReadFile(filename)
	var books []models.Booking
	json.Unmarshal(d, &books)
	var inter []interface{}
	for _, item := range books {
		inter = append(inter, item)
	}
	err := dao.Insert("Booking", "BookModel", inter...)
	if err != nil {
		fmt.Println("Insert booking error", err)
	}
}
