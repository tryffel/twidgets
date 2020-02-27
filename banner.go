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

package twidgets

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Selectable is a primitive that's able to set blur func
type Selectable interface {
	tview.Primitive
	// SetBlurFunc sets blur function that gets called upon blurring primitive
	SetBlurFunc(func(key tcell.Key))
}


// Banner combines grid layout and form-movement. To use, configure grid and add elements to it. To allow some item
// to be selected, add it to Banner.Selectable. Order in this array is same as with selections. Only buttons are
// supported as selectables. Most of the logic has been copied tview.Form.
type Banner struct {
	*tview.Grid

	//Selectable are all primitives that can be navigated and selected inside Banner.
	Selectable []Selectable
	selected int
	hasFocus bool

}

func NewBanner() *Banner {
	b := &Banner{
		Grid: tview.NewGrid(),
		Selectable: []Selectable{},
	}
	return b

}

// Focus. Copied from tview/form
func (b *Banner) Focus(delegate func(p tview.Primitive)) {
	if len(b.Selectable) == 0 {
		b.hasFocus = true
		return
	}
	b.hasFocus = false

	// Hand on the focus to one of our child elements.
	if b.selected< 0 || b.selected>= len(b.Selectable) {
		b.selected= 0
	}
	handler := func(key tcell.Key) {
		switch key {
		case tcell.KeyTab, tcell.KeyEnter, tcell.KeyCtrlJ:
			b.selected++
			b.Focus(delegate)
		case tcell.KeyBacktab, tcell.KeyCtrlK:
			b.selected--
			if b.selected< 0 {
				b.selected= len(b.Selectable) - 1
			}
			b.Focus(delegate)
		case tcell.KeyEscape:
			/*
			if b.cancel != nil {
				b.cancel()
			} else {
				b.focusedElement = 0
				b.Focus(delegate)
			}
			 */
		}
	}

	// We're selecting a button.
	button := b.Selectable[b.selected]
	button.SetBlurFunc(handler)
	delegate(button)
}
