package main

import (
	"gioui.org/io/key"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TextEditor struct {
	*Core
	txt *string
	KeyFocus
	ed widget.Editor
}

func NewTextEditor(core *Core, txt *string) *TextEditor {
	ed := &TextEditor{
		Core: core,
		txt:  txt,
	}
	ed.ed.SetText(*txt)
	return ed
}

func (ed *TextEditor) Edit(gtx C) { gtx.Execute(key.FocusCmd{Tag: &ed.ed}) }

func (ed *TextEditor) Update(gtx C) {
	for {
		e, ok := gtx.Event(
			key.Filter{Focus: &ed.ed, Name: key.NameReturn},
			key.Filter{Focus: &ed.ed, Name: key.NameEscape},
			key.Filter{Focus: &ed.ed, Name: "Z", Required: key.ModShortcut, Optional: key.ModShift},
		)
		if !ok {
			break
		}
		if e, ok := e.(key.Event); ok && e.State == key.Press {
			txt := ed.ed.Text()
			oldtxt := *ed.txt
			switch e.Name {
			case key.NameReturn:
				if txt == "" {
					break
				}
				if oldtxt == "" {
					ed.setText(gtx, txt)
					ed.Save()
					break
				}
				ed.Do(func() {
					ed.setText(gtx, txt)
				}, func() {
					ed.setText(gtx, oldtxt)
				})
			case key.NameEscape:
				if oldtxt != "" {
					ed.setText(gtx, oldtxt)
				} else {
					ed.Undo()
				}
			case "Z":
				ed.setText(gtx, oldtxt)
				if e.Modifiers.Contain(key.ModShift) {
					ed.Redo()
				} else {
					ed.Undo()
				}
			}
		}
	}
}

func (ed *TextEditor) setText(gtx C, txt string) {
	*ed.txt = txt
	ed.ed.SetText(txt)
	ed.Focus(gtx)
}

func (ed *TextEditor) Layout(gtx C) D {
	ed.Update(gtx)

	return ed.KeyFocus.Layout(ed.Core.Th, gtx, material.Editor(ed.Core.Th, &ed.ed, "").Layout)
}
