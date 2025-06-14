package main

type ToDo struct {
	ToDo, Done []*ToDoItem
}

type ToDoItem struct {
	Description string
}
