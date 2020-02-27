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

//ListItem is an item that can be used in ScrollList. Additional SetSelected is required
// since item doesn't receive focus on selection but still gets to change its visual style.
type ListItem interface {
	tview.Primitive
	SetSelected(selected bool)
}


//ScrollGrid is a list that can have more items than it can currently show.
// It allows user to scroll items. It also manages rows dynamically. Use Padding and ItemHeight to change
// grid size. Use Up/Down + (vim: j/k/g/G) to navigate between items and Enter to select item.
type ScrollList struct {
	*tview.Grid
	// Padding is num of rows or relative expansion, see tview.Grid.SetColumns() for usage
	Padding int
	ItemHeight int
	items  []ListItem
	selected int
	// range that is visible from array
	visibleFrom int
	visibleTo int
	rows int

	selectFunc func(int)
	blurFunc   func(key tcell.Key)
}

func (s *ScrollList) SetBlurFunc(blur func(key tcell.Key)) {
	s.blurFunc = blur
}

//NewScrollList creates new scroll grid. selectFunc is called whenever user presses Enter on some item.
//SelectFunc can be nil.
func NewScrollList(selectFunc func(index int)) *ScrollList {
	s := &ScrollList{
		Grid:            tview.NewGrid(),
		items:           make([]ListItem, 0),
		selectFunc:      selectFunc,
	}

	s.Padding = 1
	s.ItemHeight = 3
	return s
}

//AddItem appends single item
func (s *ScrollList) AddItem(i ListItem) {
	s.items = append(s.items, i)
	if len(s.items) == 0 {
		s.items[s.selected].SetSelected(true)
	}
	if len(s.items) < s.rows {
		s.updateGridItems()
	}
}

//AddItems appends multiple items
func (s *ScrollList) AddItems(i ...ListItem) {
	s.items = append(s.items, i...)
	s.updateGridItems()
}


//Clear clears list items an updates view
func (s *ScrollList) Clear() {
	s.items = make([]ListItem, 0)
	s.selected = 0
	s.Grid.Clear()
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

		if scrollDown && s.selected < len(s.items)-1  {
			s.items[s.selected].SetSelected(false)
			s.selected += 1
			s.items[s.selected].SetSelected(true)
			s.updateGridItems()
		} else if scrollUp && s.selected > 0  {
			s.items[s.selected].SetSelected(false)
			s.selected -= 1
			s.items[s.selected].SetSelected(true)
			s.updateGridItems()
		} else if pageUp {
			s.items[s.selected].SetSelected(false)
			s.selected = 0
			s.items[s.selected].SetSelected(true)
			s.updateGridItems()
		} else if pagedDown {
			s.items[s.selected].SetSelected(false)
			s.selected = len(s.items) -1
			s.items[s.selected].SetSelected(true)
			s.updateGridItems()
		}
	})
}


//update grid size after resizing widget
func (s *ScrollList) updateGrid(x, y, w, h int) {
	_, _, cw, ch := s.GetRect()

	// if no change, skip
	if cw == w && ch == h {
		return
	}

	rows := h / (s.ItemHeight + s.Padding * 2 -1)

	// tview adds 1 row of empty height
	if rows * s.ItemHeight + (rows + 2) * s.Padding  >= h {
		rows -= 1

	}

	if rows == -1 {
		s.Grid.Clear()
		s.Grid.SetRows()
		s.Grid.SetColumns(2, 40, -2)
		return
	}

	s.rows = rows
	// expand row items if needed
	if s.visibleFrom == 0 {
		s.visibleTo = s.rows -1
	} else if s.visibleTo == len(s.items)-1 {
		s.visibleFrom = s.visibleTo - s.rows +1
	}

	// init grid
	gridRow := make([]int, rows * 2 + 1)

	for i := 0; i < rows; i++ {
		gridRow[i*2] = s.Padding
		gridRow[i*2+1] = s.ItemHeight
	}
	// make last row flexible size
	gridRow[len(gridRow)-1] = -1

	// fill grid
	s.Grid.Clear()
	s.Grid.SetRows(gridRow...)
	s.Grid.SetColumns(2, 40, -2)

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
	if s.selected < s.visibleFrom {
		s.visibleFrom = s.selected
		s.visibleTo = s.selected + s.rows -1

	} else if s.selected > s.visibleTo {
		s.visibleTo = s.selected
		s.visibleFrom = s.selected - s.rows +1
	}

	if s.visibleTo < 0 {
		s.visibleTo = 0
	}
	if s.visibleFrom > len(s.items) -1 {
		s.visibleFrom = len(s.items) -1
	}
	if s.visibleFrom < 0 {
		s.visibleFrom = 0
	}
	if s.visibleTo > len(s.items) -1 {
		s.visibleTo = len(s.items) -1
	}

	s.Grid.Clear()
	for i := 0; i < s.rows; i++ {
		if i > s.visibleTo || s.visibleFrom + i > len(s.items) -1 {
			break
		}
		item := s.items[s.visibleFrom + i]
		s.Grid.AddItem(item, i * 2 +1,  1, 1, 1, 4,10, false)
	}
}


