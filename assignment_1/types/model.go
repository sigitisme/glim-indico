package types

type TodoTable struct {
	// The todo ID
	// example: 1
	ID int `json:"id"`
	// The todo title
	// example: Complete project
	Title string `json:"title"`
	// The todo description
	// example: Finish coding the backend
	Description string `json:"description"`
	IsDeleted   bool   `json:"is_deleted"`
}

func (t *TodoTable) ConstructTodo() Todo {
	return Todo{
		Description: t.Description,
		ID:          t.ID,
		Title:       t.Title,
	}
}

type TodoRequest struct {
	// The todo ID
	// The todo title
	// example: Complete project
	Title string `json:"title"  binding:"required"`
	// The todo description
	// example: Finish coding the backend
	Description string `json:"description" binding:"required"`
}

type Todo struct {
	// The todo ID
	ID int `json:"id"`
	// The todo title
	// example: Complete project
	Title string `json:"title"  binding:"required"`
	// The todo description
	// example: Finish coding the backend
	Description string `json:"description" binding:"required"`
}
