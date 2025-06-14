package main

type Undo struct {
	items []undoItem
	i     int
}

type undoItem struct {
	do, undo func()
}

func (u *Undo) Do(do, undo func()) {
	u.items = append(u.items[:u.i], undoItem{do: do, undo: undo})
	u.i++
	do()
}

func (u *Undo) Undo() {
	if u.i > 0 {
		u.i--
		u.items[u.i].undo()
	}
}

func (u *Undo) Redo() {
	if u.i < len(u.items) {
		do := u.items[u.i].do
		u.i++
		do()
	}
}
