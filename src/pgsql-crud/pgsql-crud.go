package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "root"
	dbPassword = ""
	dbName     = "postgres"
)

// Book - A book info
type Book struct {
	_id    int
	title  string
	author string
}

// SQLDB - Books DB Object / Bookshelf
type SQLDB struct {
	db    *sql.DB
	table string
}

// Create Table
func (bookshelf *SQLDB) createTable() (err error) {
	que := `
	CREATE TABLE IF NOT EXISTS "` + bookshelf.table + `"
	(
		"_id" serial NOT NULL,
		"title" character varying(255) NOT NULL,
		"author" character varying(255) NOT NULL,
		-- "created" date,
		-- created_at timestamp with time zone DEFAULT current_timestamp,
		CONSTRAINT userinfo_pkey PRIMARY KEY ("_id")
	)
	-- ) WITH (OIDS=FALSE); // Not work with CockroachDB`

	result, err := bookshelf.db.Exec(que)
	if err != nil {
		fmt.Println("Table Creation Error: ", result, err)
	}

	return
}

// Select - cRud (one _id)
func (bookshelf *SQLDB) getBook(bookID int) (Book, error) {
	result := Book{}

	rows, err := bookshelf.db.Query(`SELECT * FROM "`+bookshelf.table+`" where "_id"=$1`, bookID)
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&result._id, &result.title, &result.author)
			if err != nil {
				fmt.Println("Get Book Error: ", err)
			}
		}
	}

	return result, err
}

// Select - cRud (All)
func (bookshelf *SQLDB) allBooks() ([]Book, error) {
	books := []Book{}

	rows, err := bookshelf.db.Query(`SELECT * FROM "` + bookshelf.table + `" order by "_id"`)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			currentBook := Book{}
			err = rows.Scan(&currentBook._id, &currentBook.title, &currentBook.author)
			if err == nil {
				books = append(books, currentBook)
			} else {
				fmt.Println("Get All Books Error: ", err)
			}
		}
	} else {
		return books, err
	}

	return books, err
}

// Insert - Crud
// func (bookshelf *SQLDB) insertBook(title, author string) (int, error) {
func (bookshelf *SQLDB) insertBook(book Book) (int, error) {
	var bookID int

	err := bookshelf.db.QueryRow(
		`INSERT INTO "`+bookshelf.table+`"("title","author") VALUES($1,$2) RETURNING _id`,
		book.title, book.author).Scan(&bookID)
	if err != nil {
		return 0, err
	}

	return bookID, err
}

// Update - crUd
func (bookshelf *SQLDB) updateBook(_id int, book Book) (int, error) {
	res, err := bookshelf.db.Exec(
		`UPDATE "`+bookshelf.table+`" SET "title"=$1,"author"=$2 WHERE "_id"=$3 RETURNING "_id"`,
		book.title, book.author, _id)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsUpdated), err
}

// Delete - cruD
func (bookshelf *SQLDB) removeBook(bookID int) (int, error) {
	res, err := bookshelf.db.Exec(`DELETE FROM "`+bookshelf.table+`" WHERE "_id"=$1`, bookID)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsDeleted), nil
}

func main() {
	bookshelf := SQLDB{}

	dbinfo := fmt.Sprintf(
		"host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode='disable'",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	bookshelf.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer bookshelf.db.Close()

	// Create table
	// bookshelf.table = "books"
	bookshelf.table = "novel"
	bookshelf.createTable()

	// Insert
	newBook := Book{title: "표본실의 청개구리", author: "현진건"}
	newBookID, err := bookshelf.insertBook(newBook)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted book ID: ", newBookID)

	// Select an item
	book, _ := bookshelf.getBook(newBookID)
	fmt.Println("Inserted Book: ", book)

	fmt.Println("---- Wrong author ----")

	// Update item
	book.author = "염상섭"
	updateCount, err := bookshelf.updateBook(newBookID, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated count: ", updateCount)

	fmt.Println("---- Author corrected ----")

	// Select all
	books, _ := bookshelf.allBooks()
	fmt.Println("Rest of books: ", books)

	// Delete an item
	deleted, err := bookshelf.removeBook(newBookID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted count: ", deleted)
}
