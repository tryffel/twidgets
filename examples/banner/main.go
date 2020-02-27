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
	"github.com/rivo/tview"
	"os"
	"tryffel.net/go/twidgets"
)

var app *tview.Application

func main() {
	app = tview.NewApplication()

	banner := twidgets.NewBanner()

	btnExit := tview.NewButton("Exit")
	btnExit.SetSelectedFunc(func() {
		app.Stop()
		fmt.Print("Exit\n")
	})

	btnOne := tview.NewButton("One")
	btnTwo := tview.NewButton("Two")

	text := tview.NewTextView()
	text.SetText("This is a text view\nanother row")

	btns := []*tview.Button{btnExit, btnOne, btnTwo}
	for _, btn := range btns {
		btn.SetLabelColor(tcell.ColorBlack)
		btn.SetLabelColorActivated(tcell.ColorGray)
		btn.SetBackgroundColor(tcell.ColorWhite)
		btn.SetBackgroundColorActivated(tcell.ColorGreen)
	}

	banner.Selectable = btns

	banner.Grid.SetRows(1, 1, 1, 1)
	banner.Grid.SetColumns(6, 2, 10, -1, 10, -1, 10, -3)
	banner.Grid.SetMinSize(1, 6)

	banner.Grid.AddItem(btnExit, 0, 0, 1, 1, 1, 5, false)
	banner.Grid.AddItem(text, 0, 2, 2, 5, 1, 10, false)
	banner.Grid.AddItem(btnOne, 3, 2, 1, 1, 1, 10, false)
	banner.Grid.AddItem(btnTwo, 3, 4, 1, 1, 1, 10,false)

	app.SetRoot(banner, true)
	app.SetFocus(banner)
	app.Run()
}

func done(label string) {
	fmt.Printf("Got label %s", label)
	os.Exit(0)

}
