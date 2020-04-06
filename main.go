package main

import (
    "errors"
    "database/sql"
    log "github.com/sirupsen/logrus"
    "github.com/lib/pq"
    "github.com/abeerupadhyay/go-postgres-example/database"
)


var ErrInvalidRating = errors.New("invalid rating. allowed range is 0-5")

// Sample entity for demonstration
type Book struct {
    Id          uint64
    ISBN        string
    Title       string
    Author      string
    PublishYear string
    Rating      float32
}

func NewBook(isbn, title, author, publishYear string, rating float32) (*Book, error) {

    if rating < 0.0 || rating > 5.0 {
        return nil, ErrInvalidRating
    }

    return &Book{
        ISBN: isbn,
        Title: title,
        Author: author,
        PublishYear: publishYear,
        Rating: rating,
    }, nil
}


// Custom db errors
var ErrISBNAlreadyExists = errors.New("ISBN already exists")
var ErrBookNotFound = errors.New("book not found")

var pgLogger = log.WithField("ctx", "postgres")

// Wrapper for basic database operations - create, read, delete
type dbOperation struct {
    db *sql.DB
}

func (repo *dbOperation) Store(book *Book) error {

    var id uint64
    query := `
    insert into book (isbn, title, author, publish_year, rating)
    values ($1, $2, $3, $4, $5) returning id;
    `
    row := repo.db.QueryRow(
        query, book.ISBN, book.Title, book.Author,
        book.PublishYear, book.Rating,
    )
    err := row.Scan(&id)
    if err == nil {
        pgLogger.WithFields(log.Fields{"book_id": id}).Info("Book created")
    }
    // Check for unique contraint violation
    if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
        return ErrISBNAlreadyExists
    }

    book.Id = id
    return nil
}

func (repo *dbOperation) GetByISBN(isbn string) (*Book, error) {

    query := `
    select id, isbn, title, author, publish_year, rating from book
    where isbn = $1;
    `
    book := &Book{}
    err := repo.db.QueryRow(query, isbn).Scan(
        &book.Id, &book.ISBN, &book.Title, &book.Author,
        &book.PublishYear, &book.Rating,
    )
    if err == sql.ErrNoRows { err = ErrBookNotFound }
    return book, err
}

func (repo *dbOperation) DeleteById(id uint64) error {

    query := "delete from book where id = $1;"

    res, err := repo.db.Exec(query, id)
    if err != nil {
        return err
    }

    count, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if count == 0 {
        err = ErrBookNotFound
    } else {
        pgLogger.WithFields(log.Fields{"book_id": id}).Info("Book deleted")
    }

    return err
}


// Run migrations as a part of the initialization process
func init() {
    err := database.RunMigrations()
    if err != nil { panic(err) }
}

func main() {

    ops := &dbOperation{
        db: database.DefaultDb(),
    }

    defer database.CloseDefaultDb()

    var err error

    book1, err := NewBook("978-3-16-148410-0", "Lorem Ipsum", "Jane Doe", "2010", 2.3)
    if err != nil { panic(err) }

    err = ops.Store(book1)
    if err != nil { panic(err) }

    book2, err := NewBook("338-7-01-825510-9", "Foo Bar", "John Doe", "2014", 4.5)
    if err != nil { panic(err) }

    err = ops.Store(book2)
    if err != nil { panic(err) }

    err = ops.DeleteById(book1.Id)
    if err != nil { panic(err) }

    book3, err := ops.GetByISBN(book2.ISBN)
    if err != nil { panic(err) }

    err = ops.DeleteById(book3.Id)
    if err != nil { panic(err) }
}
