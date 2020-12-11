package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Handler structure
type Handler struct {
	DB *sql.DB
}

// Init returns a new Database handler
func Init(username, password, name, host string, port int) *Handler {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, name)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Handler{
		db,
	}
}

// EnsureConnected pings the database to validate whether it is alive or not
func (handler *Handler) EnsureConnected() error {
	if err := handler.DB.Ping(); err != nil {
		for i := 0; i <= 3; i++ {
			log.Fatal(err)
			time.Sleep(time.Duration(5) * time.Second)
			if err := handler.DB.Ping(); err == nil {
				break
			}
		}
		return err
	}
	return nil
}
