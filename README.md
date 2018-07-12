### 使用Golang和MongoDB构建微服务  

根据 [umermansoor github](https://github.com/umermansoor/microservices)的 `Python`版本的微服务改造成 `Golang`版本  
一共有4个微服务  
* Movie Service: 是关于电影的基本信息，标题、评分等
* ShowTimes Service: 关于电影上映时间的信息
* Booking Service: 关于电影的订阅的信息
* User Service: 用户的信息


#### 要求  

* [Golang](https://golang.org) 
* [mux](https://github.com/gorilla/mux)
* [MongoDB](https://www.mongodb.com/)

#### API和文档

各个服务之间相互独立，单独的路由和单独的数据库，各个服务之间的通信是通过 `HTTP JSON`,每个服务的API的返回结果也是JSON类型，可以参考 [使用Golang和MongoDB构建 RESTful API](https://github.com/coderminer/restful)，提取出各个服务之间共同的东西，独立于服务之外，供服务调用  

* 返回的结果封装  

`helper/utils.go`

```
func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

```

* 基础数据 Entity

`models/models.go`

```
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
```

* 关于数据库的封装  

`dao/db.go`,具体的请参考 [对 mgo关于MongoDB的基础操作的封装](https://github.com/coderminer/goutil)

```
func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}
...
```

#### 服务

各个服务具体的逻辑具体的参考 [使用Golang和MongoDB构建 RESTful API](https://github.com/coderminer/restful)  

* User Service(port 8000)
* Movie Service(port 8001)
* ShowTimes Service(port 8002)
* Booking Service(port 8003)

#### 服务通信

查询某个用户的订阅的电影信息时，需要先通过 `User Service` 服务查询这个用户，根据用户名通过 `Booking Service `查询用户的订阅信息，然后通过 `Movie Service`服务查询对应的电影的信息,都是通过 `HTTP` 通信  

```
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
```

