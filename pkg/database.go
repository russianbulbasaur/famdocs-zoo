package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitDBConnection() *sql.DB {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("dbname=%s host=%s password=%s port=%s user=%s sslmode=disable", dbName, dbHost, dbPassword, dbPort, dbUser)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}
