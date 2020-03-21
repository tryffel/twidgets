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

type Selection int

const (
	Selected Selection = iota
	Blurred
	Deselected
)

//ListItem is an item that can be used in ScrollList. Additional SetSelected is required
// since item doesn't receive focus on selection but still gets to change its visual style.
type ListItem interface {
	tview.Primitive
	SetSelected(selected Selection)
}

//ScrollGrid is a list that can have more items than it can currently show.
// It allows user to scroll items. It also manages rows dynamically. Use Padding and ItemHeight to change
// grid size. Use Up/Down + (vim: j/k/g/G) to navigate between items and Enter to select item.
type ScrollList struct {
	*tview.Grid
	// Padding is num of rows or relative expansion, see tview.Grid.SetColumns() for usage
	Padding    int
	ItemHeight int
	items      []ListItem
	selected   int
	// range that is visible from array
	visibleFrom int
	visibleTo   int
	rows        int
	gridRows    []int
	border      bool

	selectFunc       func(int)
	blurFunc         func(key tcell.Key)
	indexChangedFunc func(int) bool
}

func (s *ScrollList) SetBlurFunc(blur func(key tcell.Key)) {
	s.blurFunc = blur
}

//NewScrollList creates new scroll grid. selectFunc is called whenever user presses Enter on some item.
//SelectFunc can be nil.
func NewScrollList(selectFunc func(index int)) *ScrollList {
	s := &ScrollList{
		Grid:       tview.NewGrid(),
		items:      make([]ListItem, 0),
		selectFunc: selectFunc,
	}

	s.Grid.SetColumns(2, -2, -1)
	s.Padding = 1
	s.ItemHeight = 3
	s.gridRows = []int{2, -2, -1}
	return s
}

//AddItem appends single item
func (s *ScrollList) AddItem(i ListItem) {
	s.items = append(s.items, i)
	if len(s.items) == 0 {
		s.items[s.selected].SetSelected(Selected)
	}
	if len(s.items) < s.rows {
		s.updateGridItems()
	}
}

//AddItems appends multiple items
func (s *ScrollList) AddItems(i ...ListItem) {
	s.items = append(s.items, i...)
	x, y, w, h := s.GetRect()
	s.updateGrid(x, y, w, h)
	s.updateGridItems()
}

//Clear clears list items an updates view
func (s *ScrollList) Clear() {
	s.items = make([]ListItem, 0)
	s.selected = 0
	s.Grid.Clear()
}

// SetIndexChangedFunc sets a function that gets called every time list index is about to change.
// New index is passed to function. If it returns true, index changes. If false, index is being retained.
func (s *ScrollList) SetIndexChangedFunc(indexChanged func(int) bool) {
	s.indexChangedFunc = indexChanged
}

func (s *ScrollList) SetRect(x, y, w, h int) {
	s.updateGrid(x, y, w, h)
	s.Grid.SetRect(x, y, w, h)
}

func (s *ScrollList) GetSelectedIndex() int {
	return s.selected
}

func (s *ScrollList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.Grid.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		var acceptIndexChanged func(int) bool
		if s.indexChangedFunc != nil {
			acceptIndexChanged = s.indexChangedFunc
		} else {
			acceptIndexChanged = func(int) bool { return true }
		}

		key := event.Key()
		r := event.Rune()

		scrollDown := false
		scrollUp := false

		pagedDown := false
		pageUp := false

		switch key {
		case tcell.KeyDown:
			scrollDown = true
		case tcell.KeyUp:
			scrollUp = true
		case tcell.KeyEnter:
			if s.selectFunc != nil {
				s.selectFunc(s.selected)
			}
		case tcell.KeyTAB, tcell.KeyBacktab:
			if s.blurFunc != nil {
				s.blurFunc(key)
			}
		default:
			if r == 'j' {
				scrollDown = true
			} else if r == 'k' {
				scrollUp = true
			} else if r == 'g' {
				pageUp = true
			} else if r == 'G' {
				pagedDown = true
			}
		}

		newIndex := s.selected

		if scrollDown && s.selected < len(s.items)-1 {
			newIndex += 1
		} else if scrollUp && s.selected == 0 && s.blurFunc != nil {
			s.blurFunc(tcell.KeyBacktab)
		} else if scrollUp && s.selected > 0 {
			newIndex -= 1
		} else if pageUp {
			newIndex = 0
		} else if pagedDown {
			newIndex = len(s.items) - 1
		}

		if len(s.items) > 0 {
			if acceptIndexChanged(newIndex) {
				s.items[s.selected].SetSelected(Deselected)
				s.selected = newIndex
				s.items[s.selected].SetSelected(Selected)
				s.updateGridItems()
			}
		}
	})
}

// SetSelected sets active index. First item is 0. If value is out of bounds, do nothing.
func (s *ScrollList) SetSelected(index int) {
	if index < 0 || index > len(s.items)-1 {
		return
	}
	s.items[s.selected].SetSelected(Deselected)
	s.selected = index
	s.items[s.selected].SetSelected(Selected)
	s.updateGridItems()
}

//update grid size after resizing widget
func (s *ScrollList) updateGrid(x, y, w, h int) {
	if s.border {
		h -= 2
	}

	// how many rows with padding
	rows := h / (s.ItemHeight + s.Padding)
	takes := rows * (s.ItemHeight + s.Padding)
	noBottomPadding := false
	if takes+s.ItemHeight == h {
		// add one row without bottom padding
		rows += 1
		noBottomPadding = true
	}

	if rows == 0 {
		s.Grid.Clear()
		s.Grid.SetRows()
		return
	}

	s.rows = rows
	// expand row items if needed
	if s.visibleFrom == 0 {
		s.visibleTo = s.rows - 1
	} else if s.visibleTo == len(s.items)-1 {
		s.visibleFrom = s.visibleTo - s.rows + 1
	}

	// init grid
	gridRow := make([]int, rows*2)

	// fill rows
	for i := 0; i < rows; i++ {
		gridRow[i*2] = s.ItemHeight
		gridRow[i*2+1] = s.Padding
	}

	if !noBottomPadding {
		// set bottom padding flexible
		gridRow[len(gridRow)-1] = -1

	} else {

		gridRow = gridRow[:len(gridRow)-1]
	}

	// fill grid
	s.Grid.Clear()
	s.gridRows = gridRow
	s.Grid.SetRows(gridRow...)

	s.updateGridItems()
}

// update grid items after selecting new items
func (s *ScrollList) updateGridItems() {
	if len(s.items) == 0 {
		return
	}

	if s.rows == 1 {
		s.visibleFrom = s.selected
		s.visibleTo = s.selected
		s.Grid.Clear()
		s.Grid.AddItem(s.items[s.selected], 1, 1, 1, 1, 4, 10, false)
		return
	}

	// which items are visible, is selected one of them
	if s.selected < s.visibleFrom {
		s.visibleFrom = s.selected
		s.visibleTo = s.selected + s.rows - 1

	} else if s.selected > s.visibleTo {
		s.visibleTo = s.selected
		s.visibleFrom = s.selected - s.rows + 1
	}

	if s.visibleTo < 0 {
		s.visibleTo = 0
	}
	if s.visibleFrom > len(s.items)-1 {
		s.visibleFrom = len(s.items) - 1
	}
	if s.visibleFrom < 0 {
		s.visibleFrom = 0
	}
	if s.visibleTo > len(s.items)-1 {
		s.visibleTo = len(s.items) - 1
	}

	s.Grid.Clear()
	for i := 0; i < s.rows; i++ {
		if i > s.visibleTo || s.visibleFrom+i > len(s.items)-1 {
			break
		}
		item := s.items[s.visibleFrom+i]
		s.Grid.AddItem(item, i*2, 1, 1, 1, 4, 10, false)
	}
}

func (s *ScrollList) Focus(delegate func(p tview.Primitive)) {
	if len(s.items) > 0 {
		s.items[s.selected].SetSelected(Selected)
	}
}

func (s *ScrollList) Blur() {
	if len(s.items) > 0 {
		s.items[s.selected].SetSelected(Blurred)
	}
}

func (s *ScrollList) SetBorder(b bool) *tview.Box {
	s.Grid.SetBorder(b)
	s.border = true
	return s.Box
}
