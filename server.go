package main

import (
	"api/router"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// open .env file in the local directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("DB")
	ssl := os.Getenv("SSL_STATE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, ssl)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// srv := &http.Server{
	//     Handler:      r,
	//     Addr:         "127.0.0.1:8000",
	//     // Good practice: enforce timeouts for servers you create!
	//     WriteTimeout: 15 * time.Second,
	//     ReadTimeout:  15 * time.Second,
	// }

	r := router.Router(db) //passing the db reference to the router
	fmt.Println("Successfully connected to database!")
	log.Fatal(http.ListenAndServe(":8081", r))
}
