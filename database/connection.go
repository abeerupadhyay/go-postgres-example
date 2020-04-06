package database

import (
    "net/url"
    "os"
    "sync"
    "database/sql"
    log "github.com/sirupsen/logrus"
    _ "github.com/lib/pq"
)

// Logger used across the database package. This helps to filter
// out all database related logs by using ctx=postgres
var pgLogger = log.WithField("ctx", "postgres")


type pgConnectionUri string

func (uri pgConnectionUri) host() string {
    u, err := url.Parse(string(uri))
    if err != nil { pgLogger.Panic("Invalid connection uri.") }
    return u.Host
}

func (uri pgConnectionUri) dbName() string {
    u, err := url.Parse(string(uri))
    if err != nil { pgLogger.Panic("Invalid connection uri.") }
    return u.Path[1:]
}


// Singleton to manage database connection pool
type connectionPool struct {
    db  *sql.DB
}

func newConnectionPool(connUri pgConnectionUri) *connectionPool {

    pgConnLogger := pgLogger.WithFields(log.Fields{
        "database": connUri.dbName(),
        "host": connUri.host(),
    })

    db, err := sql.Open("postgres", string(connUri))
    if err != nil {
        pgConnLogger.Panicf("Unable to establish connection. Error: %v", err)
    }

    // send a ping to ensure connection is established
    err = db.Ping()
    if err != nil {
        pgConnLogger.Panicf("Unable to establish connection. Error: %v", err)
    }

    pgConnLogger.Infof("Connection established to database.")
    return &connectionPool{db}
}


var pgDefaultConnPool *connectionPool
var defaultOnce *sync.Once = new(sync.Once)

// Returns a connection pool object using the defualt connection uri
func DefaultDb() *sql.DB {

    if defaultOnce == nil {
        defaultOnce = new(sync.Once)
    }

    defaultOnce.Do(func() {
        connUri := pgConnectionUri(os.Getenv("PG_DEFAULT_CONN_URI"))
        pgDefaultConnPool = newConnectionPool(connUri)
    })
    return pgDefaultConnPool.db
}

// Closes the default database connection pool
func CloseDefaultDb() {
    if pgDefaultConnPool != nil {
        err := pgDefaultConnPool.db.Close()
        if err != nil {
            pgLogger.Warnf("Error closing default connection. Error: %v", err)
        }
        pgLogger.Infof("Default connection closed!")
        defaultOnce = nil
    }
}

// ---------------------------------------------------------------------
// Similarly you can create another pool for read replicas

// var pgReadOnlyConnPool *connectionPool
// var readOnlyOnce *sync.Once = new(sync.Once)

// func ReadOnlyDb() *sql.DB {
//     readOnlyOnce.Do(func() {
//         connURI := pgConnectionUri(os.Getenv("PG_READONLY_CONN_URI"))
//         pgReadOnlyConnPool = newConnectionPool(connUri)
//     })
//     return pgReadOnlyConnPool.db
// }
