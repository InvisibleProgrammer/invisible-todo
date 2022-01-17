package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HealthcheckResponse struct {
	Service bool
}

type CreateTodoItemRequest struct {
	Title       string
	Description string
}

type TodoItemRequest struct {
	Id          int32
	Title       string
	Description string
}

type TodoItem struct {
	Id          int32
	Title       string
	Description string
}

var todoItems []TodoItem

func health(c *gin.Context) {
	response := HealthcheckResponse{
		Service: true,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func addTodo(c *gin.Context) {
	var newTodoRequest CreateTodoItemRequest

	err := c.BindJSON(&newTodoRequest)
	if err != nil {
		return
	}

	var maxId int32 = 0
	for i := 0; i < len(todoItems); i++ {
		if todoItems[i].Id > maxId {
			maxId = todoItems[i].Id
		}
	}

	newTodoItem := TodoItem{
		Id:          maxId + 1,
		Title:       newTodoRequest.Title,
		Description: newTodoRequest.Description,
	}

	todoItems = append(todoItems, newTodoItem)
}

func updateTodo(c *gin.Context) {
	var updateTodoItemRequest TodoItemRequest

	idParam := c.Param("id")

	id64, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return
	}

	id := int32(id64)

	err = c.BindJSON(&updateTodoItemRequest)
	if err != nil {
		return
	}

	var itemIndex int32 = -1
	for i := 0; i < len(todoItems); i++ {
		if todoItems[i].Id == id {
			itemIndex = int32(i)

			break
		}
	}

	if itemIndex < 0 {
		return
	}

	updatedTodoItem := TodoItem{
		Id:          updateTodoItemRequest.Id,
		Title:       updateTodoItemRequest.Title,
		Description: updateTodoItemRequest.Description,
	}

	todoItems[itemIndex] = updatedTodoItem
}

func deleteTodo(c *gin.Context) {
	idParam := c.Param("id")

	id64, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return
	}

	id := int32(id64)

	var itemIndex int32 = -1
	for i := 0; i < len(todoItems); i++ {
		if todoItems[i].Id == id {
			itemIndex = int32(i)

			break
		}
	}

	if itemIndex >= 0 {
		removeIndex(todoItems, int(itemIndex))
	}
}

func removeIndex(s []TodoItem, index int) []TodoItem {
	return append(s[:index], s[index+1:]...)
}

func listTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoItems)
}

func main() {

	todoItems = []TodoItem{}

	router := gin.Default()
	router.GET("/diag/health", health)

	router.GET("/todo", listTodos)
	router.POST("/todo", addTodo)
	router.PUT("/todo/:id", updateTodo)
	router.DELETE("/todo/:id", deleteTodo)

	router.Run("localhost:8080")
}
