# go-postgres-example  [![Build Status](https://travis-ci.com/abeerupadhyay/go-postgres-example.svg?branch=master)](https://travis-ci.com/abeerupadhyay/go-postgres-example)  [![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/abeerupadhyay/go-postgres-example/?ref=repository-badge)

Go application to demonstrate integration with postgresql

To run the script locally, set `PG_DEFAULT_CONN_URI=postgres://myuser:pass123@localhost:5432/mydb?sslmode=disable` and `PG_MIGRATIONS_FILE_PATH=/path/to/project/database/migrations` in your environment variables and run `go run main.go`. If the values were set correctly, you must see the following output

```shell
INFO[0000] Migrations applied successfully               ctx=postgres
INFO[0000] Connection established to database.           ctx=postgres database=localdb host="localhost:5432"
INFO[0000] Book created                                  book_id=1 ctx=postgres
INFO[0000] Book created                                  book_id=2 ctx=postgres
INFO[0000] Book deleted                                  book_id=1 ctx=postgres
INFO[0000] Book deleted                                  book_id=2 ctx=postgres
INFO[0000] Default connection closed!                    ctx=postgres
```
