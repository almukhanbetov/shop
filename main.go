package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден, читаем системные переменные")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL не указан")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("PostgreSQL не отвечает:", err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend работает",
		})
	})

	router.GET("/persons", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, age FROM persons ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		persons := []Person{}

		for rows.Next() {
			var person Person

			err := rows.Scan(&person.ID, &person.Name, &person.Age)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			persons = append(persons, person)
		}

		c.JSON(http.StatusOK, persons)
	})

	log.Println("Backend запущен: http://localhost:" + port)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}