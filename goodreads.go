package goodreads

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	id    string
	key   string
	limit int

	books   = Books{}
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

type Books map[string]Book

type Book struct {
	Id       string `xml:"id"`
	Title    string `xml:"title"`
	Link     string `xml:"link"`
	ImageURL string `xml:"image_url"`
	NumPages string `xml:"num_pages"`
}

// PUBLIC

func SetConfig(i, k string, l int) {
	id = i
	key = k
	limit = l
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
