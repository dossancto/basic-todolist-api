package main

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
  "os"

	"context"
	"fmt"
	"net/http"

	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TodoItem struct {
	Id              int    `json:"id"`
	TodoName        string `json:"todo_name"`
	TodoDescription string `json:"todo_description"`
	IsChecked       bool   `json:"is_checked"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var dbCtx, cancel = context.WithTimeout(context.Background(), 60*5*time.Second)

var database *sql.DB

func getTodoById(ctx *gin.Context) {
	var todoItem = TodoItem{}

	todoId := ctx.Param("id")

	rows, err := database.Query("SELECT * from tbTodos where id = '" + todoId + "'")
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&todoItem.Id, &todoItem.TodoName, &todoItem.TodoDescription, &todoItem.IsChecked)
		checkError(err)
	}

	err = rows.Err()
	checkError(err)

	ctx.IndentedJSON(http.StatusOK, todoItem)
}

func getAllTodoNames(ctx *gin.Context) {
	var allNames []string

	rows, err := database.Query("SELECT todoName from tbTodos;")

	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var todoName string

		err := rows.Scan(&todoName)

		checkError(err)

		allNames = append(allNames, todoName)
	}

	err = rows.Err()
	checkError(err)
	ctx.IndentedJSON(http.StatusOK, allNames)

}

func all(ctx *gin.Context) {
	var allTodos = []TodoItem{}

	rows, err := database.Query("SELECT * from tbTodos;")

	checkError(err)
	defer rows.Close()

	for rows.Next() {
		todoItem := TodoItem{}

		err := rows.Scan(&todoItem.Id, &todoItem.TodoName, &todoItem.TodoDescription, &todoItem.IsChecked)

		checkError(err)

		allTodos = append(allTodos, todoItem)
	}

	err = rows.Err()
	checkError(err)
	ctx.IndentedJSON(http.StatusOK, allTodos)
}

func searchTodoByName(ctx *gin.Context) {
	searchString := ctx.Param("name")
	var allTodos = []TodoItem{}

	queryString := "SELECT * from tbTodos where todoName like '%" + searchString + "%';"

	rows, err := database.Query(queryString)

	checkError(err)
	defer rows.Close()

	for rows.Next() {
		todoItem := TodoItem{}

		err := rows.Scan(&todoItem.Id, &todoItem.TodoName, &todoItem.TodoDescription, &todoItem.IsChecked)

		checkError(err)

		allTodos = append(allTodos, todoItem)
	}

	err = rows.Err()
	checkError(err)
	ctx.IndentedJSON(http.StatusOK, allTodos)
}


func main() {
	var dotenvErr = godotenv.Load()

	var host = os.Getenv("DBHOST")
	var databaseName = os.Getenv("DBNAME")
	var user = os.Getenv("DBUSER")
	var password = os.Getenv("DBPASSWORD")

  checkError(dotenvErr)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowNativePasswords=true", user, password, host, databaseName)

	var err error

	database, err = sql.Open("mysql", connectionString)
	checkError(err)
	defer database.Close()

	err = database.Ping()
	checkError(err)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/all", all)
	router.GET("/api/names", getAllTodoNames)

	router.GET("/api/:id", getTodoById)

	router.GET("/api/search/:name", searchTodoByName)

	// router.POST("/api", addMapa)

	router.Run()
}
