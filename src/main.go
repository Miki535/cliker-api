package main

import (
	"database/sql"
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

type DataToPutInDB struct {
	Username string
	Record   int
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
		var dataToPut DataToPutInDB

		if err := c.ShouldBindJSON(&dataToPut); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad JSON request",
			})
			return
		}

		insertUser(db, dataToPut.Username, dataToPut.Record)

		c.JSON(200, gin.H{
			"result:": "all good!",
		})
	})

	router.GET("/getDatabase", func(c *gin.Context) {
		users, err := selectUsersFromDB(db)
		if err != nil {
			log.Fatal("Error while selecting users data from DB: ", err)
		}

		c.JSON(200, gin.H{
			"Data: ": users,
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

func insertUser(db *sql.DB, username string, record int) {
	insertUserSQL := `INSERT INTO users(username, record) VALUES (?, ?)`
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	_, err = statement.Exec(username, record)
	if err != nil {
		log.Println("Error insetring user: ", err.Error())
		return
	}
}

func selectUsersFromDB(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, record FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Record); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
