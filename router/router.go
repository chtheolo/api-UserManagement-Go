package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB /*declare it global, so we can use it in every HandleFunc.*/

type user struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	Bday      string `json:"bday"`
}

func returnUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT * FROM "users"`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	users := make([]user, 0)
	for rows.Next() {
		u := user{}
		err := rows.Scan(&u.Firstname, &u.Lastname, &u.Address, &u.Bday, &u.Id) //order matters
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}
	usersJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(usersJson)
	if err != nil {
		panic(err)
	}

	for _, usr := range users {
		fmt.Printf("%s %s %s %s %s", usr.Id, usr.Firstname, usr.Lastname, usr.Address, usr.Bday)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// get form values
	u := user{}
	u.Firstname = r.FormValue("firstname")
	u.Lastname = r.FormValue("lastname")
	u.Address = r.FormValue("address")
	u.Bday = r.FormValue("bday")

	// validate the form values
	if u.Firstname == "" || u.Lastname == "" || u.Address == "" || u.Bday == "" {
		fmt.Println(r.FormValue("firstname"))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (FIRSTNAME, LASTNAME, ADDRESS, BDAY) VALUES ($1, $2, $3, $4)", u.Firstname, u.Lastname, u.Address, u.Bday)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	key := vars["id"]

	row := db.QueryRow(`SELECT * FROM "users" WHERE id = $1`, key)
	u := user{}
	err := row.Scan(&u.Firstname, &u.Lastname, &u.Address, &u.Bday, &u.Id)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userJson)
	if err != nil {
		panic(err)
	}

	// fmt.Fprintf(w, "%s, %s, %s, %s, %s", u.Id, u.Firstname, u.Lastname, u.Address, u.Bday)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	key := vars["id"]

	u := user{}
	u.Firstname = r.FormValue("firstname")
	u.Lastname = r.FormValue("lastname")
	u.Address = r.FormValue("address")
	u.Bday = r.FormValue("bday")

	// validate the form values
	if u.Firstname == "" || u.Lastname == "" || u.Address == "" || u.Bday == "" {
		fmt.Println(r.FormValue("firstname"))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// insert new values
	_, err := db.Exec("UPDATE users SET firstname = $1, lastname = $2, address = $3, bday = $4 WHERE id = $5", u.Firstname, u.Lastname, u.Address, u.Bday, key)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	key := vars["id"]

	_, err := db.Exec("DELETE FROM users WHERE id=$1", key)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

/*export Router*/
func Router(database *sql.DB) *mux.Router {
	fmt.Println("Router enabled ...")
	db = database
	Router := mux.NewRouter().StrictSlash(true)

	Router.HandleFunc("/users", returnUsers).Methods("GET")
	Router.HandleFunc("/users", createUser).Methods("POST")
	Router.HandleFunc("/users/{id}", returnSingleUser).Methods("GET")
	Router.HandleFunc("/users/{id}", editUser).Methods("PUT")
	Router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	return Router
}
