package main

import (
	"database/sql"
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := conDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/info", func(c *gin.Context) {
		data, err := getInfo(db)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, data)
	})

	r.Run(":5050")
}

func conDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/world")
	if err != nil {
		log.Fatal("connect to database fail:", err)
	}
	return db
}

func getInfo(db *sql.DB) ([]map[string]interface{}, error) {
	info, err := db.Query("SELECT id, content, own FROM info")
	if err != nil {
		return nil, err
	}
	defer info.Close()

	var results []map[string]interface{}
	for info.Next() {
		var id int
		var content string
		var own string
		if err := info.Scan(&id, &content, &own); err != nil {
			return nil, err
		}
		data := map[string]interface{}{
			"id":      id,
			"content": content,
			"own":     own,
		}
		results = append(results, data)
	}

	if err := info.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
