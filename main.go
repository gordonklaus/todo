package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		if err := Main(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func Main() error {
	var w app.Window
	w.Option(app.Title("To Do"))
	once := sync.OnceFunc(func() { w.Perform(system.ActionCenter) })

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	name := filepath.Join(home, "Documents", "TODO")
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	core := NewCore(name)

	if err := core.Load(); err != nil {
		return err
	}

	ed := NewToDoEditor(core, core.ToDo)

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			once()

			ops.Reset()
			gtx := app.NewContext(&ops, e)

			// Disable tab navigation globally.
			for ok := true; ok; _, ok = gtx.Event(key.Filter{Name: key.NameTab, Optional: key.ModShift}) {
			}

			paint.Fill(gtx.Ops, core.Th.Bg)
			ed.Layout(gtx)

			e.Frame(&ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}
