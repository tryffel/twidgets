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

package twidgets

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type NavBarColors struct {
	Background            tcell.Color
	BackgroundFocus       tcell.Color
	ButtonBackground      tcell.Color
	ButtonBackgroundFocus tcell.Color
	Text                  tcell.Color
	TextFocus             tcell.Color
	Shortcut              tcell.Color
	ShortcutFocus         tcell.Color
}

/* NavBar implements navigation bar with multiple buttons
 */
type NavBar struct {
	grid      *tview.Grid
	buttons   []*tview.Button
	btnKeys   []tcell.Key
	btnLabels []string
	doneFunc  func(label string)
	colors    *NavBarColors

	// Which button is active
	btnActiveIndex int
	hasFocus       bool
}

func (n *NavBar) Draw(screen tcell.Screen) {
	n.grid.Draw(screen)
}

func (n *NavBar) GetRect() (int, int, int, int) {
	return n.grid.GetRect()
}

func (n *NavBar) SetRect(x, y, width, height int) {
	n.grid.SetRect(x, y, width, height)
}

func (n *NavBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		lastBtn := n.btnActiveIndex

		key := event.Key()
		if key == tcell.KeyRight {
			n.btnActiveIndex = min(len(n.buttons)-1, n.btnActiveIndex+1)
		} else if key == tcell.KeyLeft {
			n.btnActiveIndex = max(0, n.btnActiveIndex-1)
		}

		if lastBtn != n.btnActiveIndex {
			n.buttons[lastBtn].Blur()
			n.buttons[n.btnActiveIndex].Focus(nil)
		}

		if key == tcell.KeyEnter {
			n.callDone(n.btnLabels[n.btnActiveIndex])
		}

		for i, v := range n.btnKeys {
			if key == v {
				n.callDone(n.btnLabels[i])
				break
			}
		}
		n.grid.InputHandler()(event, setFocus)
	}
}

func (n *NavBar) Focus(delegate func(p tview.Primitive)) {
	n.grid.Focus(delegate)
	n.hasFocus = true
	n.buttons[n.btnActiveIndex].Focus(nil)
}

func (n *NavBar) Blur() {
	n.btnActiveIndex = 0
	n.hasFocus = false
	n.grid.Blur()
}

func (n *NavBar) GetFocusable() tview.Focusable {
	return n.grid.GetFocusable()
}

func NewNavBar(colors *NavBarColors, doneFunc func(label string)) *NavBar {
	nav := &NavBar{
		grid:      tview.NewGrid(),
		buttons:   []*tview.Button{},
		btnKeys:   []tcell.Key{},
		btnLabels: []string{},
		doneFunc:  doneFunc,
		colors:    colors,
	}

	nav.grid.SetBorders(false)
	nav.grid.SetBorder(false)
	nav.grid.SetBackgroundColor(colors.Background)
	nav.grid.SetRows(-1)

	return nav
}

func (n *NavBar) AddButton(button *tview.Button, key tcell.Key) {
	n.buttons = append(n.buttons, button)
	n.btnKeys = append(n.btnKeys, key)
	n.btnLabels = append(n.btnLabels, button.GetLabel())
	button.SetBorder(false)
	button.SetSelectedFunc(wrapKeyFunc(button.GetLabel(), n.doneFunc))
	button.SetBackgroundColor(n.colors.ButtonBackground)
	button.SetBackgroundColorActivated(n.colors.ButtonBackgroundFocus)
	button.SetLabelColor(n.colors.Text)
	button.SetLabelColorActivated(n.colors.TextFocus)

	hex := n.colors.Shortcut.Hex()

	label := fmt.Sprintf("[#%02x]%s[-] %s", hex, tcell.KeyNames[key], button.GetLabel())
	button.SetLabel(label)

	count := len(n.buttons)
	n.grid.AddItem(button, 0, 2*count, 1, 1, 1, 5, false)

	widths := make([]int, len(n.buttons)*2+1)
	spaceWidth := -1
	for i := 0; i < len(n.buttons)+1; i++ {
		widths[i*2] = -2
		if i > 0 {
			widths[i*2-1] = spaceWidth
		}
	}

	n.grid.SetColumns(widths...)
}

func wrapKeyFunc(label string, doneFunc func(label string)) func() {
	return func() {
		doneFunc(label)
	}
}

func (n *NavBar) callDone(label string) {
	if n.doneFunc != nil {
		n.doneFunc(label)
	}
}
