/*
 * Copyright 2020 Tero Vierimaa
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
	"log"
	"tryffel.net/go/twidgets"
)

type Item struct {
	*cview.TextView
	text string
}

func (i *Item) SetSelected(selected twidgets.Selection) {
	if selected == twidgets.Selected {
		i.SetBorderAttributes(tcell.AttrBold)
		i.SetBorderColor(tcell.ColorBlue)
		i.SetTextColor(tcell.ColorYellow)
	} else if selected == twidgets.Deselected {
		i.SetBorderAttributes(tcell.AttrNone)
		i.SetBorderColor(tcell.ColorWhite)
		i.SetTextColor(tcell.ColorGray)
	}
}

func NewItem(text string) *Item {
	i := &Item{
		TextView: cview.NewTextView(),
		text:     text,
	}

	i.SetBorder(true)
	i.SetText(text)
	return i
}

func main() {
	app := cview.NewApplication()
	list := twidgets.NewScrollList(nil)
	list.ItemHeight = 2

	items := make([]*Item, 0)

	for i := 0; i < 10; i++ {
		item := NewItem(fmt.Sprintf("%d. item\nrow2", i))
		item.SetBorder(false)
		item.SetBorderPadding(0, 0, 1, 1)
		list.AddItem(item)
		items = append(items, item)
	}

	// alter item text when it's selected / deselected
	indexChangedFunc := func(next int) bool {
		current := list.GetSelectedIndex()
		items[current].SetText(items[current].text)

		items[next].SetText("  " + items[next].text)
		return true
	}

	list.SetIndexChangedFunc(indexChangedFunc)

	list.AddContextItem("Option a", 0, func(n int) {})
	list.AddContextItem("Option b", 0, func(n int) {})
	list.AddContextItem("Option c", 0, func(n int) {})
	list.AddContextItem("Option d", 0, func(n int) {})
	list.AddContextItem("Option e", 0, func(n int) {})

	app.SetRoot(list, true)
	app.Run()
}

func printSelect(index int) {
	log.Println(fmt.Sprintf("selected index %d", index))
}
