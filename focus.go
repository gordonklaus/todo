package main

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type KeyFocus struct {
	_ int // because pointers to empty structs may not be unique
}

func (f *KeyFocus) Focus(gtx C)        { gtx.Execute(key.FocusCmd{Tag: f}) }
func (f *KeyFocus) Focused(gtx C) bool { return gtx.Focused(f) }

func (f *KeyFocus) Event(gtx C, e *key.Event, required, optional key.Modifiers, names ...key.Name) bool {
	filters := make([]event.Filter, 1+len(names))
	filters[0] = key.FocusFilter{Target: f}
	for i, name := range names {
		filters[1+i] = key.Filter{Focus: f, Required: required, Optional: optional, Name: name}
	}
	for {
		ev, ok := gtx.Event(filters...)
		if !ok {
			return false
		}
		if ev, ok := ev.(key.Event); ok && ev.State == key.Press {
			*e = ev
			return true
		}
	}
}

func (f *KeyFocus) Layout(th *material.Theme, gtx C, w layout.Widget) D {
	event.Op(gtx.Ops, f)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if f.Focused(gtx) {
				m := gtx.Dp(4)
				defer op.Offset(image.Pt(-m, -m)).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, mulAlpha(th.Fg, 64),
					clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min.Add(image.Pt(2*m, 2*m))}, 2*m).Op(gtx.Ops))
			}
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}
