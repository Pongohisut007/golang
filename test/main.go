package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/world")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err

}

func getUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "User")

	row, err := db.Query("SELECT NAME FROM nongao")
	if err != nil {
		http.Error(w, "ไม่สามารถดึงข้อมูลได้", 500)
	}
	defer row.Close()

	var names []string

	for row.Next() {
		var name string
		err = row.Scan(&name)
		if err != nil {
			http.Error(w, "ไม่สามารถอ่านข้อมูลได้", 500)
		}
		names = append(names, name)
	}
	for _, n := range names {
		fmt.Fprintf(w, "%s\n", n)
	}

}

type Person struct {
	Name string `json:name`
	Age  int    `json:age`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allow", 405)
		return
	}

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "json error", 400)
		return
	}

	_, err = db.Exec("INSERT INTO nongao (name, age) VALUES (?, ?)", p.Name, p.Age)
	if err != nil {
		http.Error(w, "internal server", 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	resp := map[string]string{"status": "success", "message": "addPerson success"}
	json.NewEncoder(w).Encode(resp)

}

type CreateA struct {
	Name     string `json:"name" validate:"required"`
	Password int    `json:"password" validate:"required,min=1"`
	Sex      string `json:"sex" validate:"required,oneof=male female"`
}

var validate = validator.New()

func createpassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allow", 405)
		return
	}

	var c CreateA
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, "json inva lid", 400)
		return
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		http.Error(w, "Invalid form format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO user (name, password, sex) VALUES (?, ?, ?)", c.Name, c.Password, c.Sex)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	resp := map[string]string{
		"status:": "success",
		"message": "create user password success"}
	json.NewEncoder(w).Encode(resp)

}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s Url: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("fail to connect db")
	}
	fmt.Println("connect to db success")
	defer db.Close()

	http.Handle("/get", logMiddleware(http.HandlerFunc(getUsers)))
	http.Handle("/adduser", logMiddleware(http.HandlerFunc(addUser)))
	http.Handle("/create", logMiddleware(http.HandlerFunc(createpassword)))

	s := &http.Server{
		Addr: ":8000",
	}
	log.Fatal(s.ListenAndServe())

}
