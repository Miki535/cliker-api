package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Username string
	Record   int
}

type ClickerBasicStructure struct {
	Record   int    `json:"record"`
	Nickname string `json:"nickname"`
}

// just for testing before real database
type TestingGetRequest struct {
	Names   string `json:"names"`
	Records int    `json:"records"`
}

func main() {
	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal("Fatal error while open DB: ", err)
	}
	defer db.Close()

	createTable(db)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Api is working!",
		})
	})
	// router for put info in "ClickerBasicStructure":) kurwa!
	router.POST("/postInformation", func(c *gin.Context) {
		var clickerBaseST ClickerBasicStructure

		if err := c.ShouldBindJSON(&clickerBaseST); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad JSON request",
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "All good!",
		})
	})

	router.GET("/getDatabase", func(c *gin.Context) {
		var testingGT TestingGetRequest

		testingGT.Names = "John deer"
		testingGT.Records = 1785983

		c.JSON(200, gin.H{
			"Names: ":   testingGT.Names,
			"Records: ": testingGT.Records,
		})
	})

	router.Run()
}

func createTable(db *sql.DB) {
	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"username" TEXT NOT NULL UNIQUE,
		"record" INTEGER
	);`

	log.Println("Creating users table...")
	statement, err := db.Prepare(createUsersTableSQL)
	if err != nil {
		log.Fatal("Fatal error while creating table: ", err.Error())
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("users table created")
}

func insertUser(db *sql.DB, username string, age int) {
	insertUserSQL := `INSERT INTO users(username, record) VALUES (?, ?)`
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	_, err = statement.Exec(username, age)
	if err != nil {
		log.Println("Error insetring user: ", err.Error())
		return
	}
}

func selectUsers(db *sql.DB) {
	row, err := db.Query("SELECT * FROM users ORDER BY id")
	if err != nil {
		log.Fatal("Error while selecting from DB", err)
	}
	defer row.Close()

	var users []User
	for row.Next() {
		var user User

		err := row.Scan(&user.ID, &user.Username, &user.Record)
		if err != nil {
			log.Fatal("Error while scanning row: ", err)
		}
		users = append(users, user)
	}
	err = row.Err()
	if err != nil {
		log.Fatal("Error: type 'row.Err()'", err)
	}

	fmt.Println("Users in db + this is testing lines of code. so chilll:)")
	for _, user := range users {
		fmt.Printf("ID: %d, Username: %s, Age: %d\n", user.ID, user.Username, user.Record)
	}
}
