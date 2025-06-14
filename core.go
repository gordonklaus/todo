package main

import (
	"errors"
	"image/color"
	"io/fs"
	"log"
	"os"

	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/todo/internal/todo"
)

type Core struct {
	filename string
	ToDo     *ToDo
	Th       *material.Theme
	undo     Undo
}

func NewCore(filename string) *Core {
	th := material.NewTheme()
	th.Fg = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	th.Bg = color.NRGBA{R: 16, G: 16, B: 16, A: 255}
	th.ContrastFg = color.NRGBA{R: 16, G: 16, B: 16, A: 0}
	th.ContrastBg = color.NRGBA{R: 64, G: 192, B: 192, A: 255}

	return &Core{
		filename: filename,
		Th:       th,
	}
}

func (c *Core) Do(do, undo func()) { c.undo.Do(do, undo); c.Save() }
func (c *Core) Undo()              { c.undo.Undo(); c.Save() }
func (c *Core) Redo()              { c.undo.Redo(); c.Save() }

func (c *Core) Load() error {
	f, err := os.Open(c.filename)
	if errors.Is(err, fs.ErrNotExist) {
		c.ToDo = &ToDo{}
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	var todo todo.ToDo
	if err := bits.Read(f, &todo); err != nil {
		return err
	}

	t := &ToDo{
		ToDo: make([]*ToDoItem, len(todo.ToDo)),
		Done: make([]*ToDoItem, len(todo.Done)),
	}
	for i, it := range todo.ToDo {
		t.ToDo[i] = &ToDoItem{Description: it.Description}
	}
	for i, it := range todo.Done {
		t.Done[i] = &ToDoItem{Description: it.Description}
	}
	c.ToDo = t

	return nil
}

func (c *Core) Save() {
	f, err := os.Create(c.filename)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	tt := &todo.ToDo{
		ToDo: make([]todo.Item, len(c.ToDo.ToDo)),
		Done: make([]todo.Item, len(c.ToDo.Done)),
	}
	for i, it := range c.ToDo.ToDo {
		tt.ToDo[i] = todo.Item{Description: it.Description}
	}
	for i, it := range c.ToDo.Done {
		tt.Done[i] = todo.Item{Description: it.Description}
	}

	if err := bits.Write(f, tt); err != nil {
		log.Print(err)
		return
	}
}
