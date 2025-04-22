package main

import (
	"fmt"
	// "todo_app_go/handlers"
	// "todo_app_go/models"
	"todo_app_go/storage"
)

func main() {
	fmt.Println("todo_app")

	storage.InitDB()
}
