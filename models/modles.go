package models

type User struct {
	Id   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Movie struct {
	Id       string  `bson:"_id" json:"id"`
	Title    string  `bson:"title" json:"title"`
	Rating   float32 `bson:"rating" json:"rating"`
	Director string  `bson:"director" json:"director"`
}

type ShowTimes struct {
	Id     string   `bson:"_id" json:"id"`
	Date   string   `bson:"date" json:"date"`
	Movies []string `bson:"movies" json:"movies"`
}

type Booking struct {
	Id    string     `bson:"_id" json:"id"`
	Name  string     `bson:"name" json:"name"`
	Books []BookInfo `bson:"books" json:"books"`
}

type BookInfo struct {
	Date   string   `bson:"date" json:"date"`
	Movies []string `bson:"movies" json:"movies"`
}

type Result struct {
	Name  string       `json:"name"`
	Books []ResultInfo `json:"books"`
}

type ResultInfo struct {
	Date   string  `json:"date"`
	Movies []Movie `json:"movies"`
}
