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
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	author := query.Get("author")
	title := query.Get("title")
	price := query.Get("price")

	if len(price) != 0 || len(author) != 0 || len(title) != 0 {
		json.NewEncoder(w).Encode(filterByQuery(title, author, price))
	} else {
		json.NewEncoder(w).Encode(books)
	}
}

func filterByQuery(title, author, price string) []Book {
	var filteredBooks []Book

	for i := 0; i < len(books); i++ {

		var currentBook Book = books[i]
		if (len(author) == 0 || checkStrings(author, currentBook.Author)) &&
			(len(title) == 0 || checkStrings(title, currentBook.Title)) &&
			(len(price) == 0 || checkNumber(price, currentBook.Price)) {
			filteredBooks = append(filteredBooks, currentBook)
		}

	}
	return filteredBooks
}

func checkStrings(searchedString, checkWith string) bool {
	return strings.Contains(strings.ToLower(checkWith), strings.ToLower(searchedString))
}

func checkNumber(searchedNumber, checkWith string) bool {
	if number, err := strconv.Atoi(searchedNumber); err == nil {
		if bookPrice, err := strconv.Atoi(checkWith); err == nil {
			return bookPrice <= number
		}
	}
	return false
}

//Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get Params
	//Loop through books to find id
	for i := 0; i < len(books); i++ {
		if books[i].Id == params["id"] {
			json.NewEncoder(w).Encode(books[i])
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

//Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook Book
	_ = json.NewDecoder(r.Body).Decode(&newBook)
	IdNum += 1
	newBook.Id = strconv.Itoa(IdNum)
	books = append(books, newBook)
	json.NewEncoder(w).Encode(newBook)
	writeToJson()
}

//Update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i := 0; i < len(books); i++ {
		if books[i].Id == params["id"] {
			//Slicing through and removing that index
			books = append(books[:i], books[i+1:]...)
			var secondHalf []Book
			secondHalf = append(secondHalf, books[i:]...)
			var updatedbook Book
			_ = json.NewDecoder(r.Body).Decode(&updatedbook)
			updatedbook.Id = params["id"]
			//Slicing through and adding updated book to removed index
			books = append(append(books[:i], updatedbook), secondHalf...)
			json.NewEncoder(w).Encode(updatedbook)
			writeToJson()
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

//Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.Id == params["id"] {
			//Slicing through and removing that index
			books = append(books[:index], books[index+1:]...)
			writeToJson()
			break
		}
	}
	json.NewEncoder(w).Encode(books)
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

	//To get the max id present
	for _, item := range books {
		if maxid, err := strconv.Atoi(item.Id); err == nil {
			if maxid > IdNum {
				IdNum = maxid
			}

		}
	}

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
