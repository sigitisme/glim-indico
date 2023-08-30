package handlers

import (
	"assignment_1/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	data = make([]types.TodoTable, 0)
)

// @Summary Get all todos
// @Description Get a list of all todos
// @Tags todos
// @ID get-todos
// @Produce json
// @Success 200 {array} []types.Todo
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/v1/todos [get]
func GetTodos(c *gin.Context) {
	todos := []types.Todo{}

	for _, v := range data {
		if v.IsDeleted {
			continue
		}

		todo := v.ConstructTodo()

		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

// @Summary Get todo by id
// @Description Get a todo by id
// @Tags todos
// @ID get-todo-by-id
// @Produce json
// @Param        id   path      int  true  "Todo ID"
// @Success 200 {array} types.Todo
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router /api/v1/todos/{id} [get]
func GetTodoByID(c *gin.Context) {
	id := c.Param("id")

	idInt, httpError := CheckId(id)

	if httpError != nil {
		c.JSON(httpError.Status, httpError)
		return
	}

	value := data[idInt-1]

	if value.IsDeleted {
		c.JSON(http.StatusNotFound, HTTPError{
			Status:  http.StatusNotFound,
			Message: "Invalid ID",
		})
		return
	}

	c.JSON(http.StatusOK, value.ConstructTodo())
}

// @Summary Create a todo
// @Description Create a todo
// @Tags todos
// @ID create-todo
// @Produce json
// @Param request body types.TodoRequest.request true "request body"
// @Success 200 {array} types.Todo
// @Router /api/v1/todos [post]
func CreateTodo(c *gin.Context) {
	todoReq := types.TodoRequest{}

	if err := c.ShouldBindJSON(&todoReq); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Invalid Request Body",
		})
		return
	}

	count := len(data)

	todo := types.TodoTable{
		Description: todoReq.Description,
		Title:       todoReq.Title,
		ID:          count + 1,
	}

	data = append(data, todo)

	c.JSON(200, todo.ConstructTodo())
}

// @Summary Update a todo by id
// @Description Update a todo by id
// @Tags todos
// @ID update-todo-by-id
// @Produce json
// @Param        id   path      int  true  "Todo ID"
// @Success 200 {array} types.Todo
// @Router /api/v1/todos/{id} [put]
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")

	idInt, httpError := CheckId(id)

	if httpError != nil {
		c.JSON(httpError.Status, httpError)
		return
	}

	todoReq := types.TodoRequest{}

	if err := c.ShouldBindJSON(&todoReq); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Invalid Request Body",
		})
		return
	}

	value := data[idInt-1]

	if value.IsDeleted {
		c.JSON(http.StatusNotFound, HTTPError{
			Status:  http.StatusNotFound,
			Message: "Invalid ID",
		})
		return
	}

	todo := types.TodoTable{
		Description: todoReq.Description,
		Title:       todoReq.Title,
		ID:          idInt,
	}

	data[idInt-1] = todo

	response := todo.ConstructTodo()

	c.JSON(http.StatusOK, response)
}

// @Summary Delete a todo
// @Description Delete a todo
// @Tags todos
// @ID delete-todo-by-id
// @Produce json
// @Param        id   path      int  true  "Todo ID"
// @Success 200 {array} types.Todo
// @Router /api/v1/todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	idInt, httpError := CheckId(id)

	if httpError != nil {
		c.JSON(httpError.Status, httpError)
		return
	}

	data[idInt-1].IsDeleted = true

	c.Status(http.StatusOK)
}
