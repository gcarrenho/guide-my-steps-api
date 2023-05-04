package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "man1234"
	dbname   = "DB_1"
)

type postgreSQLConf struct {
	dbName string // The database that you want to access
	dbUser string // The database account that you want to access
	dbHost string // The endpoint of the DB instance that you want to access
	dbPort int    // The port number used for connecting to your DB instance
	//region string // The AWS Region where the DB instance is running
}

func NewPostgreSQLConf() postgreSQLConf {
	return postgreSQLConf{
		dbName: os.Getenv("DBNAME"),
		dbUser: os.Getenv("MYSQLUSER"),
		dbHost: os.Getenv("DBHOST"),
		dbPort: port,
	}
}

func (postgresql *postgreSQLConf) InitPostgreSQLDB() (*sql.DB, error) {
	password := os.Getenv("POSTGRESQLPASSWORD")
	//connStr := "postgres://postgres:password@localhost/DB_1?sslmode=disable"
	//dsn := fmt.Sprintf("postgresql://%s:%s@tcp(%s:%d)/%s?parseTime=true",
	//postgresql.dbUser, password, mysql.dbHost, mysql.dbPort, mysql.dbName)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresql.dbHost, postgresql.dbPort, postgresql.dbUser, password, postgresql.dbName)

	// Connect to database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
