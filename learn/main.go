package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Person struct ใช้สำหรับรับ/ส่งข้อมูล JSON
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// handler /hi ดึงชื่อจากตาราง nongao
func sayhi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi")

	rows, err := db.Query("SELECT NAME FROM nongao")
	if err != nil {
		http.Error(w, "can't pull data", 500)
		return
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			http.Error(w, "can't read data", 500)
			return
		}
		names = append(names, name)
	}
	for _, n := range names {
		fmt.Fprintln(w, n)
	}
}

// handler รับ JSON POST และเพิ่มข้อมูลลงฐานข้อมูล
func addPersonToDBHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	_, err = db.Exec("INSERT INTO nongao (name, age) VALUES (?, ?)", p.Name, p.Age)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"status": "success", "mess age": "เพิ่มข้อมูลเรียบร้อยแล้ว"}
	json.NewEncoder(w).Encode(resp)
}

// middleware log request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// connectDB เชื่อมต่อฐานข้อมูล MySQL
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/world")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("เชื่อมต่อ db ไม่ได้:", err)
	}
	log.Println("เชื่อมต่อ db สำเร็จ")

	http.Handle("/hi", loggingMiddleware(http.HandlerFunc(sayhi)))
	http.Handle("/addpersonjson", loggingMiddleware(http.HandlerFunc(addPersonToDBHandler)))

	log.Println("server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
