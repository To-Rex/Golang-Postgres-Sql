package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user	 = "postgres"
	password = "0000"
	dbname   = "postgres"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	Promocode string `json:"promocode"`
	Status   string `json:"status"`
	Roles    string `json:"roles"`
}

func passwordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func main() {
	r := gin.Default()
	r.POST("/register", register)
	r.Run(":8080")
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}


func register(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := connectDB()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (username VARCHAR(255) PRIMARY KEY, email VARCHAR(255), password VARCHAR(255), name VARCHAR(255), surname VARCHAR(255), age INT,  phone VARCHAR(255), promocode VARCHAR(255), status VARCHAR(255), roles VARCHAR(255))")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Promocode = "test"
	user.Status = "active"
	user.Roles = "admin"
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is empty"})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is empty"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is empty"})
		return
	}
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is empty"})
		return
	}
	if user.Surname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "surname is empty"})
		return
	}
	if user.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is empty"})
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).Scan(&username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exist"})
		return
	}

	var email string
	err = db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email).Scan(&email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exist"})
		return
	}

	_, err = db.Exec("INSERT INTO users (username, email, password, name, surname, age, phone, promocode, status, roles) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", user.Username, user.Email, passwordHash(user.Password), user.Name, user.Surname, user.Age, user.Phone, user.Promocode, user.Status, user.Roles)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error insert data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}