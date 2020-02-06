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
	"github.com/rivo/tview"
	"log"
	"tryffel.net/go/twidgets"
)

type Item struct {
	*tview.TextView
}

func (i *Item) SetSelected(selected bool) {
	if selected {
		i.SetBorderAttributes(tcell.AttrBold)
		i.SetBorderColor(tcell.ColorBlue)
	} else {
		i.SetBorderAttributes(tcell.AttrNone)
		i.SetBorderColor(tcell.ColorWhite)
	}
}

func NewItem(text string) *Item {
	i := &Item{
		TextView: tview.NewTextView(),
	}

	i.SetBorder(true)
	i.SetText(text)
	return i
}


func main() {
	app := tview.NewApplication()
	list := twidgets.NewScrollList(printSelect)
	list.ItemHeight = 5

	for i := 0; i < 10; i++ {
		item := NewItem(fmt.Sprintf("%d. item\nrow2\nrow3", i))
		item.SetBorderPadding(0, 0, 1, 1)
		list.AddItem(item)
	}
	app.SetRoot(list, true)
	app.Run()
}

func printSelect(index int) {
	log.Println(fmt.Sprintf("selected index %d", index))
}
