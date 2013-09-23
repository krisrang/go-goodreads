package goodreads

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	id     string
	key    string
	secret string

	apiRoot = "http://www.goodreads.com/user/show/"
)

type Response struct {
	User User `xml:"user"`
}

type User struct {
	ID            string `xml:"id"`
	Name          string `xml:"name"`
	About         string `xml:"about"`
	Link          string `xml:"link"`
	ImageURL      string `xml:"image_url"`
	SmallImageURL string `xml:"small_image_url"`
	Location      string `xml:"location"`
	LastActive    string `xml:"last_active"`
	Reviews       int    `xml:"reviews_count"`
}

func (u User) ReviewsLink() string {
	return "http://www.goodreads.com/review/list/" + u.ID + "?sort=review&view=reviews"
}

// PUBLIC

func SetConfig(i, k, s string) {
	id = i
	key = k
	secret = s
}

func GetUser() *User {
	uri := apiRoot + id + ".xml?key=" + key
	response := &Response{}
	getData(uri, response)
	return &response.User
}

// PRIVATE

func getData(uri string, i interface{}) {
	data := getRequest(uri)
	xmlUnmarshal(data, i)
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func xmlUnmarshal(b []byte, i interface{}) {
	err := xml.Unmarshal(b, i)
	if err != nil {
		log.Fatal(err)
	}
}
