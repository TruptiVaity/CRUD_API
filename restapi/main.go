package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

//Book
type Book struct {
	Id     string `json:"id"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  string `json:"price"`
}

//Init books var as a slice Book
var books []Book
var IdNum int

//Get All Books
func getBooks(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	query := request.URL.Query()
	author := query.Get("author")
	title := query.Get("title")
	price := query.Get("price")

	if len(price) != 0 || len(author) != 0 || len(title) != 0 {
		json.NewEncoder(responseWriter).Encode(filterByQuery(title, author, price))
	} else {
		json.NewEncoder(responseWriter).Encode(books)
	}
}

func filterByQuery(title, author, price string) []Book {
	var filteredBooks []Book

	for i := 0; i < len(books); i++ {

		var currentBook Book = books[i]
		if filterByAuthor(author, currentBook) &&
			filterByTitle(title, currentBook) &&
			filterByPrice(price, currentBook) {
			filteredBooks = append(filteredBooks, currentBook)
		}

	}
	return filteredBooks
}

func filterByAuthor(searchedAuthor string, currentBook Book) bool {
	return len(searchedAuthor) == 0 || checkStrings(searchedAuthor, currentBook.Author)
}

func filterByTitle(searchedTitle string, currentBook Book) bool {
	return len(searchedTitle) == 0 || checkStrings(searchedTitle, currentBook.Title)
}

func checkStrings(searchedString, checkWith string) bool {
	return strings.Contains(strings.ToLower(checkWith), strings.ToLower(searchedString))
}

func filterByPrice(searchedPrice string, currentBook Book) bool {
	if len(searchedPrice) != 0 {
		if number, err := strconv.Atoi(searchedPrice); err == nil {
			if bookPrice, err := strconv.Atoi(currentBook.Price); err == nil {
				return bookPrice <= number
			}
		}
	}
	return false
}

//Get Single Book
func getBook(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) //Get Params
	//Loop through books to find id
	for i := 0; i < len(books); i++ {
		if books[i].Id == params["id"] {
			json.NewEncoder(responseWriter).Encode(books[i])
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(&Book{})
}

//Create a new book
func createBook(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var newBook Book
	_ = json.NewDecoder(request.Body).Decode(&newBook)
	newBook.Id = strconv.Itoa(int(idGenerator()))
	books = append(books, newBook)
	json.NewEncoder(responseWriter).Encode(newBook)
	writeToJson()
}

//Update a book
func updateBook(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for i := 0; i < len(books); i++ {
		if books[i].Id == params["id"] {
			books[i].Author = params["author"]
			books[i].Title = params["title"]
			books[i].Price = params["price"]
			_ = json.NewDecoder(request.Body).Decode(&books[i])
			json.NewEncoder(responseWriter).Encode(books[i])
			writeToJson()
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(books)
}

//Delete a book
func deleteBook(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range books {
		if item.Id == params["id"] {
			//Slicing through and removing that index
			books = append(books[:index], books[index+1:]...)
			writeToJson()
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(books)
}

func main() {
	//Read json file
	jsonFile, err := os.Open("books.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//Add values to struct
	json.Unmarshal(byteValue, &books)
	defer jsonFile.Close()

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/api/books", getBooks).Methods("GET")
	myRouter.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	myRouter.HandleFunc("/api/books", createBook).Methods("POST")
	myRouter.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	myRouter.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func writeToJson() {
	file, _ := json.MarshalIndent(books, "", "")
	ioutil.WriteFile("books.json", file, 0644)
}

func idGenerator() int {

	guid := xid.New()
	return int(guid.Counter())

}
