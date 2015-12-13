package goodreads

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	apiRoot = "http://www.goodreads.com/"
)

type Response struct {
	User    User     `xml:"user"`
	Book    Book     `xml:"book"`
	Reviews []Review `xml:"reviews>review"`
}

type User struct {
	ID            string       `xml:"id"`
	Name          string       `xml:"name"`
	About         string       `xml:"about"`
	Link          string       `xml:"link"`
	ImageURL      string       `xml:"image_url"`
	SmallImageURL string       `xml:"small_image_url"`
	Location      string       `xml:"location"`
	LastActive    string       `xml:"last_active"`
	ReviewCount   int          `xml:"reviews_count"`
	Statuses      []UserStatus `xml:"user_statuses>user_status"`
	Shelves       []Shelf      `xml:"user_shelves>user_shelf"`
	LastRead      []Review
}

func (u User) ReadingShelf() Shelf {
	for _, shelf := range u.Shelves {
		if shelf.Name == "currently-reading" {
			return shelf
		}
	}

	return Shelf{}
}

func (u User) ReadShelf() Shelf {
	for _, shelf := range u.Shelves {
		if shelf.Name == "read" {
			return shelf
		}
	}

	return Shelf{}
}

func (u User) ToReadShelf() Shelf {
	for _, shelf := range u.Shelves {
		if shelf.Name == "to-read" {
			return shelf
		}
	}

	return Shelf{}
}

type Shelf struct {
	ID        string `xml:"id"`
	BookCount string `xml:"book_count"`
	Name      string `xml:"name"`
}

type UserStatus struct {
	Page    int    `xml:"page"`
	Percent int    `xml:"percent"`
	Updated string `xml:"updated_at"`
	Book    Book   `xml:"book"`
}

func (u UserStatus) UpdatedRelative() string {
	return relativeDate(u.Updated)
}

type Book struct {
	ID       string   `xml:"id"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	ImageURL string   `xml:"image_url"`
	NumPages string   `xml:"num_pages"`
	Format   string   `xml:"format"`
	Authors  []Author `xml:"authors>author"`
	ISBN     string   `xml:"isbn"`
}

func (b Book) Author() Author {
	return b.Authors[0]
}

type Author struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
	Link string `xml:"link"`
}

type Review struct {
	Book   Book   `xml:"book"`
	Rating int    `xml:"rating"`
	ReadAt string `xml:"read_at"`
	Link   string `xml:"link"`
}

func (r Review) FullStars() []bool {
	return make([]bool, r.Rating)
}

func (r Review) EmptyStars() []bool {
	return make([]bool, 5-r.Rating)
}

func (r Review) ReadAtShort() string {
	date, err := parseDate(r.ReadAt)
	if err != nil {
		return ""
	}

	return (string)(date.Format("2 Jan 2006"))
}

func (r Review) ReadAtRelative() string {
	return relativeDate(r.ReadAt)
}

// PUBLIC

func GetUser(id, key string, limit int) *User {
	uri := apiRoot + "user/show/" + id + ".xml?key=" + key
	response := &Response{}
	getData(uri, response)

	for i := range response.User.Statuses {
		status := &response.User.Statuses[i]
		bookid := status.Book.ID
		book := GetBook(bookid, key)
		status.Book = book
	}

	if len(response.User.Statuses) >= limit {
		response.User.Statuses = response.User.Statuses[:limit]
	} else {
		remaining := limit - len(response.User.Statuses)
		response.User.LastRead = GetLastRead(id, key, remaining)
	}

	return &response.User
}

func GetBook(id, key string) Book {
	uri := apiRoot + "book/show/" + id + ".xml?key=" + key
	response := &Response{}
	getData(uri, response)

	return response.Book
}

func GetLastRead(id, key string, limit int) []Review {
	l := strconv.Itoa(limit)
	uri := apiRoot + "review/list/" + id + ".xml?key=" + key + "&v=2&shelf=read&sort=date_read&order=d&per_page=" + l

	response := &Response{}
	getData(uri, response)

	return response.Reviews
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

func parseDate(s string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		date, err = time.Parse(time.RubyDate, s)
		if err != nil {
			return time.Time{}, err
		}
	}

	return date, nil
}

func relativeDate(d string) string {
	date, err := parseDate(d)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	s := time.Now().Sub(date)

	days := int(s / (24 * time.Hour))
	if days > 1 {
		return fmt.Sprintf("%v days ago", days)
	} else if days == 1 {
		return fmt.Sprintf("%v day ago", days)
	}

	hours := int(s / time.Hour)
	if hours > 1 {
		return fmt.Sprintf("%v hours ago", hours)
	}

	minutes := int(s / time.Minute)
	if minutes > 2 {
		return fmt.Sprintf("%v minutes ago", minutes)
	} else {
		return "Just now"
	}
}
