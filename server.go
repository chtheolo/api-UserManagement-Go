package main

import (
	"api/router"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "$k0p3l0$"
		dbname   = "users"
		ssl      = "disable"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, ssl)
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

	r := router.Router(db)
	fmt.Println("Successfully connected to database!")
	log.Fatal(http.ListenAndServe(":8081", r))
}
