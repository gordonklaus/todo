package main

import (
	"slices"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ToDoEditor struct {
	*Core
	todo, done *ToDoListEditor
}

func NewToDoEditor(core *Core, todo *ToDo) *ToDoEditor {
	ed := &ToDoEditor{
		Core: core,
	}
	ed.todo = NewToDoListEditor(core, ed, &todo.ToDo)
	ed.done = NewToDoListEditor(core, ed, &todo.Done)
	return ed
}

func (ed *ToDoEditor) Update(gtx C) {
	ed.todo.Update(gtx)
	ed.done.Update(gtx)

	for {
		e, ok := gtx.Event(key.Filter{Name: "Z", Required: key.ModShortcut, Optional: key.ModShift})
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.Event:
			if e.State == key.Press {
				switch e.Name {
				case "Z":
					if e.Modifiers.Contain(key.ModShift) {
						ed.Redo()
					} else {
						ed.Undo()
					}
				}
			}
		}
	}
}

func (ed *ToDoEditor) Layout(gtx C) D {
	ed.Update(gtx)

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, ed.todo.Layout),
		layout.Flexed(1, ed.done.Layout),
	)
}

type ToDoListEditor struct {
	*Core
	parent      *ToDoEditor
	items       *[]*ToDoItem
	list        widget.List
	editors     []*ToDoItemEditor
	placeholder KeyFocus
}

func NewToDoListEditor(core *Core, parent *ToDoEditor, items *[]*ToDoItem) *ToDoListEditor {
	ed := &ToDoListEditor{
		Core:    core,
		parent:  parent,
		items:   items,
		list:    widget.List{List: layout.List{Axis: layout.Vertical}},
		editors: make([]*ToDoItemEditor, len(*items)),
	}
	for i := range *items {
		ed.editors[i] = NewToDoItemEditor(core, (*ed.items)[i])
	}
	return ed
}

func (ed *ToDoListEditor) new(gtx C, i int, after bool) {
	if after {
		i++
	}
	i = max(0, min(i, len(ed.editors)))

	ited := NewToDoItemEditor(ed.Core, &ToDoItem{})
	ed.Do(
		func() { ed.insert(gtx, i, ited) },
		func() { ed.delete(gtx, i, after) },
	)
}

func (ed *ToDoListEditor) insert(gtx C, i int, ited *ToDoItemEditor) {
	*ed.items = slices.Insert(*ed.items, i, ited.item)
	ed.editors = slices.Insert(ed.editors, i, ited)
	if ited.item.Description == "" {
		ited.Edit(gtx)
	} else {
		ited.Focus(gtx)
	}
}

func (ed *ToDoListEditor) delete(gtx C, i int, back bool) {
	*ed.items = slices.Delete(*ed.items, i, i+1)
	ed.editors = slices.Delete(ed.editors, i, i+1)
	if back {
		i--
	}
	i = max(0, min(i, len(*ed.items)-1))
	if i >= 0 && i < len(*ed.items) {
		ed.editors[i].Focus(gtx)
	} else {
		ed.placeholder.Focus(gtx)
	}
}

func (ed *ToDoListEditor) Focus(gtx C) {
	if len(ed.editors) > 0 {
		ed.editors[0].Focus(gtx)
	} else {
		ed.placeholder.Focus(gtx)
	}
}

func (ed *ToDoListEditor) Update(gtx C) {
	if gtx.Focused(nil) {
		ed.Focus(gtx)
	}

	for i, ited := range slices.Clone(ed.editors) {
		var e key.Event
		switch {
		case ited.Event(gtx, &e, 0, key.ModShortcut|key.ModShift, key.NameReturn):
			if e.Modifiers == key.ModShortcut {
				ited.Edit(gtx)
			} else {
				ed.new(gtx, i, e.Modifiers == 0)
			}
		case ited.Event(gtx, &e, 0, key.ModShift, key.NameLeftArrow, key.NameRightArrow):
			led1 := ed.parent.todo
			led2 := ed.parent.done
			if e.Name == key.NameLeftArrow {
				led1, led2 = led2, led1
			}
			if ed == led1 {
				if e.Modifiers == key.ModShift {
					ed.Do(func() {
						*led1.items = slices.Delete(*led1.items, i, i+1)
						led1.editors = slices.Delete(led1.editors, i, i+1)
						*led2.items = slices.Insert(*led2.items, 0, ited.item)
						led2.editors = slices.Insert(led2.editors, 0, ited)
						ited.Focus(gtx)
					}, func() {
						*led1.items = slices.Insert(*led1.items, i, ited.item)
						led1.editors = slices.Insert(led1.editors, i, ited)
						*led2.items = slices.Delete(*led2.items, 0, 1)
						led2.editors = slices.Delete(led2.editors, 0, 1)
						ited.Focus(gtx)
					})
				} else {
					led2.Focus(gtx)
				}
			}
		case ited.Event(gtx, &e, 0, key.ModShift, key.NameUpArrow, key.NameDownArrow):
			if e.Modifiers == key.ModShift {
				if e.Name == key.NameUpArrow {
					i--
				}
				j := i + 1
				if i >= 0 && j < len(*ed.items) {
					ed.Do(func() {
						(*ed.items)[i], (*ed.items)[j] = (*ed.items)[j], (*ed.items)[i]
						ed.editors[i], ed.editors[j] = ed.editors[j], ed.editors[i]
						ited.Focus(gtx)
					}, func() {
						(*ed.items)[i], (*ed.items)[j] = (*ed.items)[j], (*ed.items)[i]
						ed.editors[i], ed.editors[j] = ed.editors[j], ed.editors[i]
						ited.Focus(gtx)
					})
				}
				break
			}
			if e.Name == key.NameUpArrow && i > 0 {
				ed.editors[i-1].Focus(gtx)
			} else if e.Name == key.NameDownArrow && i < len(ed.editors)-1 {
				ed.editors[i+1].Focus(gtx)
			}
		case ited.Event(gtx, &e, 0, 0, key.NameDeleteBackward, key.NameDeleteForward):
			ed.Do(
				func() { ed.delete(gtx, i, e.Name == key.NameDeleteBackward) },
				func() { ed.insert(gtx, i, ited) },
			)
		}

		ited.Update(gtx)
	}

	var e key.Event
	switch {
	case ed.placeholder.Event(gtx, &e, 0, 0, key.NameReturn):
		ed.new(gtx, 0, false)
	case ed.placeholder.Event(gtx, &e, 0, 0, key.NameLeftArrow):
		if ed == ed.parent.done {
			ed.parent.todo.Focus(gtx)
		}
	case ed.placeholder.Event(gtx, &e, 0, 0, key.NameRightArrow):
		if ed == ed.parent.todo {
			ed.parent.done.Focus(gtx)
		}
	}
}

func (ed *ToDoListEditor) Layout(gtx C) D {
	ed.Update(gtx)

	return layout.UniformInset(8).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lbl := material.H5(ed.Core.Th, "To Do")
				if ed == ed.parent.done {
					lbl.Text = "Done"
				}
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				if len(ed.editors) > 0 {
					return material.List(ed.Core.Th, &ed.list).Layout(gtx, len(ed.editors), func(gtx C, index int) D {
						return layout.UniformInset(4).Layout(gtx, ed.editors[index].Layout)
					})
				}
				return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
					return ed.placeholder.Layout(ed.Core.Th, gtx, layout.Spacer{Width: 4, Height: gtx.Metric.SpToDp(ed.Core.Th.TextSize * 4 / 3)}.Layout)
				})
			}),
		)
	})
}

type ToDoItemEditor struct {
	*Core
	item *ToDoItem
	*TextEditor
}

func NewToDoItemEditor(core *Core, item *ToDoItem) *ToDoItemEditor {
	return &ToDoItemEditor{
		Core:       core,
		item:       item,
		TextEditor: NewTextEditor(core, &item.Description),
	}
}
