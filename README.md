Description:

An API using Golang providing standard CRUD operations on a resource.
Resource: Books Data
Parameters of a Book: ISBN, Author, Title, Price

API consists of following criteria:
- An endpoint to create a new item: createBook
- An endpoint to read an item: getBook
- An endpoint to read all items: getBooks

	o Set this endpoint up 3 different query parameters used to filter the books: 
        
    1. Filter books with Author name		
    2. Filter books with Title 
    3. Filter books with price less than the entered value

- An endpoint to update an item: updateBook.
- An endpoint to delete an item: deleteBook.

This API collects the data from a JSON file "books.json". (For testing purposes, ID's from the books.json file can be used to verify "getBook, updateBook and deleteBook" functions)

Create new book entry: Takes the input. A unique ID is created using the Counter method from Package xid which is a globally unique id generator library. 

Update a book: Find the match for the requested Id and update its parameters and send response.

Delete a book: Removes the entry from the data.

Read All items:
	If query is present, filter the books by query. Else return all the items present in the datastore.

FilterByQuery:
1. Query with name "author" checks for the authors present in the datastore.
2. Query with name "title" checks for book titles.
3. Query with name "price" checks for all the books with less than or equal to the value entered by the user. 
4. Using more than one query at a time separated by "&" helps in narrowing your search results.
	
Instructions to Run the project:
1. Open command prompt
2. Go the home folder of your project (here "restapi")
3. Run "go run main.go"
4. Download or open "Postman" or any similar API development friendly application to check POST, PUT, DELETE operations.
6. Read all items: http://localhost:8081/api/books (GET)
7. Get single item : http://localhost:8081/api/books/{id} (GET)
8. Create a new item: http://localhost:8081/api/books/ (POST)
9. Update an item: http://localhost:8081/api/books/{id} (PUT)
10. Delete an item: http://localhost:8081/api/books/{id} (DELETE)
11. Filter Query(GET): 
	
    a. Check if the data contains books by entered author: e.g. http://localhost:8081/api/books?author="author"
	
    b. Check if the data contains books with entered title: e.g. http://localhost:8081/api/books?author="title"
	
    c. List all the books with prices less than the entered price: e.g. http://localhost:8081/api/books?price="150"
	
    d. You can apply more than one filter at a time: e.g. http://localhost:8081/api/books?author="author"&price="150&title="title"

