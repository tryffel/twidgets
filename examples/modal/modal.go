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
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
	"tryffel.net/go/twidgets"
)

// Some modal to show
type Modal struct {
	*tview.TextView
	doneFunc func()
}

func (m *Modal) SetDoneFunc(doneFunc func()) {
	m.doneFunc = doneFunc
}

func (m *Modal) SetVisible(visible bool) {
}

//Catch enter and escape
func (m *Modal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()

		if key == tcell.KeyEnter || key == tcell.KeyEscape {
			m.doneFunc()
		}
	}
}

// Demonstrate both Modal and ModalLayout
func main() {
	app := tview.NewApplication()

	layout := twidgets.NewModalLayout()

	modal := Modal{
		TextView: tview.NewTextView(),
	}

	modal.SetBorder(true)
	modal.SetText("A Modal. \nPress enter or escape to close this.")
	text := tview.NewTextView()
	text.SetText("Some ordinary text. \nModal opens in 1 second")
	text.SetBorder(true)

	layout.Grid().AddItem(text, 0, 0, 10, 10, 10,10,false)

	//Close modal
	close := func() {
		app.QueueUpdateDraw(func() {
			layout.RemoveModal(&modal)
			modal.Blur()
			app.SetFocus(text)
		})
	}

	//Open modal
	open := func() {
		time.Sleep(1 * time.Second)
		app.QueueUpdateDraw(func() {
			text.Blur()
			//layout.AddDynamicModal(&modal, twidgets.ModalSizeMedium/)
			layout.AddFixedModal(&modal, 10, 10, twidgets.ModalSizeLarge)

			app.SetFocus(&modal)
		})
	}

	modal.SetDoneFunc(close)

	app.SetRoot(layout, true)
	go open()
	app.Run()
}



