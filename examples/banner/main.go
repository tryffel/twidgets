/*
 *   Copyright 2019 Tero Vierimaa
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
	"os"
	"tryffel.net/go/twidgets"
)

// we need button that implements twidgets.Selectable
type button struct {
	*cview.Button
}

func (b *button) SetBlurFunc(blur func(key tcell.Key)) {
	b.Button.SetBlurFunc(blur)
}

func newButton(label string) *button {
	b := &button{
		Button: cview.NewButton(label),
	}
	return b
}

var app *cview.Application

func main() {
	app = cview.NewApplication()

	banner := twidgets.NewBanner()

	btnExit := newButton("Exit")
	btnExit.SetSelectedFunc(func() {
		app.Stop()
		fmt.Print("Exit\n")
	})

	btnOne := newButton("One")
	btnTwo := newButton("Two")

	text := cview.NewTextView()
	text.SetText("This is a text view\nanother row")

	btns := []*button{btnExit, btnOne, btnTwo}
	selectables := []twidgets.Selectable{btnExit, btnOne, btnTwo}
	for _, btn := range btns {
		btn.SetLabelColor(tcell.ColorBlack)
		btn.SetLabelColorActivated(tcell.ColorGray)
		btn.SetBackgroundColor(tcell.ColorWhite)
		btn.SetBackgroundColorActivated(tcell.ColorGreen)
	}

	banner.Selectable = selectables
	banner.Grid.SetRows(1, 1, 1, 1)
	banner.Grid.SetColumns(6, 2, 10, -1, 10, -1, 10, -3)
	banner.Grid.SetMinSize(1, 6)

	banner.Grid.AddItem(btnExit, 0, 0, 1, 1, 1, 5, false)
	banner.Grid.AddItem(text, 0, 2, 2, 5, 1, 10, false)
	banner.Grid.AddItem(btnOne, 3, 2, 1, 1, 1, 10, false)
	banner.Grid.AddItem(btnTwo, 3, 4, 1, 1, 1, 10, false)

	app.SetRoot(banner, true)

	app.SetRoot(banner, true)
	app.EnableMouse(true)
	app.SetFocus(banner)
	app.Run()
}

func done(label string) {
	fmt.Printf("Got label %s", label)
	os.Exit(0)

}
