package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func conDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(0.0.0.0:3306)/world")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

func main() {
	db, err := conDB()
	if err != nil {
		fmt.Println("Database connected error")
	}
	fmt.Println("Database connected successfully")
	defer db.Close()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hi",
		})
	})

	r.Run(":5000")
}
