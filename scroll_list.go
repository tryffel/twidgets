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

type ListItem interface {
	tview.Primitive
	SetSelected(selected bool)
}


//ScrollGrid is a grid that can have more items than it can currently show.
// It allows user to scroll items. It also manages columns and rows dynamically.
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
}

//NewScrollList creates new scroll grid. ItemHeight/weight sets single item height
func NewScrollList() *ScrollList {
	s := &ScrollList{
		Grid:            tview.NewGrid(),
		items:           make([]ListItem, 0),
	}

	s.Padding = 1
	s.ItemHeight = 3
	return s
}

func (s *ScrollList) AddItem(i ListItem) {
	s.items = append(s.items, i)
}


func (s *ScrollList) Clear() {
	s.items = make([]ListItem, 0)
	s.selected = 0
	s.Grid.Clear()
}

func (s *ScrollList) SetRect(x, y, w, h int) {
	s.updateGrid(x, y, w, h)
	s.Grid.SetRect(x, y, w, h)
}

func (s *ScrollList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.Grid.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()

		switch key {
		case tcell.KeyDown:
			if s.selected < len(s.items)-1 {
				s.items[s.selected].SetSelected(false)
				s.selected += 1
				s.items[s.selected].SetSelected(true)
				s.updateGridItems()
			}
		case tcell.KeyUp:
			if s.selected > 0 {
				s.items[s.selected].SetSelected(false)
				s.selected -= 1
				s.items[s.selected].SetSelected(true)
				s.updateGridItems()
			}
		default:
			return
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

	rows := h / s.ItemHeight
	// tview adds 1 row of empty height
	if rows * s.ItemHeight + (rows + 2) * s.Padding + 1 >= h {
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
		s.visibleTo = s.rows
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


