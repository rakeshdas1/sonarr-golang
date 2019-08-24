package main

import (
	"github.com/jroimartin/gocui"
)

type List struct {
	*gocui.View
	title       string
	items       []interface{}
	currPageIdx int
	ordered     bool
}

func CreateList(v *gocui.View, ordered bool) *List {
	list := &List{}
	list.View = v
	list.SelBgColor = gocui.ColorBlack
	list.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	list.AutoScroll = true
	list.ordered = ordered
	return list
}
func (l *List) SetItems(data []interface{}) error {
	l.items = data
	return l.Draw()
}
