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
	"gitlab.com/tslocum/cview"
)

const (
	arrowUp   = "▲"
	arrowDown = "▼"
)

// SortType is a direction that can be sorted with
type Sort int

const (
	// Sort ascending
	SortAsc Sort = iota
	// Sort descending
	SortDesc
)

// Table extends cview.Table with some helpers for managing rows.
// In addition it provides sorting capabilities.
type Table struct {
	*cview.Table

	columns          []string
	columnWidths     []int
	columnExpansions []int
	showIndex        bool
	sortCol          int
	sortType         Sort

	sortFunc    func(col string, sort Sort)
	addCellFunc func(cell *cview.TableCell, header bool, col int)
}

// NewTable creates new table instance
func NewTable() *Table {
	t := &Table{
		Table: cview.NewTable(),
	}

	t.Table.SetFixed(1, 100)
	t.Table.SetSelectable(true, false)
	t.sortCol = 0
	t.sortType = SortAsc
	t.SetCellSimple(0, 0, "#")

	t.SetFixed(1, 10)
	return t
}

// SetSortFunc sets sorting function that gets called whenever user calls sorting some column
func (t *Table) SetSortFunc(sortFunc func(column string, sort Sort)) *Table {
	t.sortFunc = sortFunc
	return t
}

// SetAddCellFunc add function callback that gets called every time a new cell is added with flag of whether
// the cell is in header row. Use this to modify e.g. style of the cell when it gets added to table.
func (t *Table) SetAddCellFunc(cellFunc func(cell *cview.TableCell, header bool, col int)) *Table {
	t.addCellFunc = cellFunc
	return t
}

// SetShowIndex configure whether first column in table is item index. If set, first item is in index 1.
// Changing this does not update existing data. Thus data needs to be cleared and rows added again for changes
// to take effect
func (t *Table) SetShowIndex(index bool) {
	t.showIndex = index
}

// SetColumnWidths sets each columns maximum width. If index is included as first row,
// it must be included in here.
func (t *Table) SetColumnWidths(widths []int) {
	t.columnWidths = widths
}

// SetColumnExpansions sets how each column will expand / shrink when changing terminal size.
// If index is included as first row, it must be included in here.
func (t *Table) SetColumnExpansions(expansions []int) {
	t.columnExpansions = expansions
}

// Clear clears the content of the table. If headers==true, remove headers as well
func (t *Table) Clear(headers bool) *Table {
	if headers {
		t.Table.Clear()
	} else {
		count := t.Table.GetColumnCount()
		cells := make([]*cview.TableCell, count)
		for i := 0; i < count; i++ {
			cells[i] = t.Table.GetCell(0, i)
		}

		t.Table.Clear()
		for i := 0; i < count; i++ {
			t.Table.SetCell(0, i, cells[i])
		}
	}
	t.SetOffset(1, 0)
	return t
}

//AddRow adds single row to table
func (t *Table) AddRow(index int, content ...string) *Table {
	count := len(content)

	cells := make([]*cview.TableCell, count, count+1)
	for i := 0; i < len(content); i++ {
		cells[i] = cview.NewTableCell(content[i])
	}

	if t.showIndex {
		cells = append([]*cview.TableCell{cview.NewTableCell(fmt.Sprint(index + 1))}, cells...)
	}

	for i := 0; i < len(cells); i++ {
		if len(t.columnWidths) >= i && t.columnWidths != nil {
			cells[i].SetMaxWidth(t.columnWidths[i])
		}
		if len(t.columnExpansions) >= i && t.columnExpansions != nil {
			cells[i].SetExpansion(t.columnExpansions[i])
		}

		if t.addCellFunc != nil {
			t.addCellFunc(cells[i], false, index+1)
		}
		t.Table.SetCell(index+1, i, cells[i])
	}
	return t
}

// SetSort sets default sort column and type
func (t *Table) SetSort(column int, sort Sort) *Table {
	if t.showIndex && column == 0 {
		t.sortCol = 1
	} else {
		t.sortCol = column
	}
	t.sortType = sort
	t.updateSort()
	return t
}

// SetColumns set column header names. This will clear the table
func (t *Table) SetColumns(columns []string) *Table {
	t.Clear(true)
	if t.showIndex {
		columns = append([]string{"#"}, columns...)
		if len(columns) >= 2 {
			t.sortCol = 1
			t.sortType = SortAsc
		}
	} else {
		if len(columns) >= 1 {
			t.sortCol = 0
			t.sortType = SortAsc
		}
	}
	for i := 0; i < len(columns); i++ {
		cell := cview.NewTableCell(columns[i])
		if t.addCellFunc != nil {
			t.addCellFunc(cell, true, 0)
		}
		t.Table.SetCell(0, i, cell)
	}
	t.columns = columns
	return t
}

//Inputhandler handles header row inputs
func (t *Table) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
		enableHeader := false
		key := event.Key()
		if t.sortFunc != nil {
			row, _ := t.Table.GetSelection()
			if row == 1 && key == tcell.KeyUp {
				enableHeader = true
				t.Table.SetSelectable(true, true)
				t.Table.Select(0, t.sortCol)
			} else if row == 0 && key == tcell.KeyDown {
				t.Table.SetSelectable(true, false)
			}
			if key == tcell.KeyEnter && row == 0 && t.sortFunc != nil {
				t.updateSort()
			}
		}
		// User might move to first/last row, catch when user moves to 0 row and select
		// 1st row instead. This is only if user moved with other key than key up
		row, _ := t.Table.GetSelection()
		atHeader := row == 0
		t.Table.InputHandler()(event, setFocus)
		row, _ = t.Table.GetSelection()
		if row == 0 && !atHeader && !enableHeader {
			t.Table.Select(1, 0)
			t.Table.SetSelectable(true, false)
		} else if enableHeader {
			t.Table.Select(0, t.sortCol)
			t.Table.SetSelectable(true, true)
		}
	}
}

//update sort and call sortFunc if there is one
func (t *Table) updateSort() {
	_, col := t.GetSelection()
	if col == 0 && t.showIndex {
		//Refuse to sort by index
		return
	}
	cell := t.GetCell(0, t.sortCol)
	if t.sortCol == col {
		if t.sortType == SortAsc {
			t.sortType = SortDesc
			cell.SetText(fmt.Sprintf("%s %s", t.columns[col], arrowUp))
		} else {
			t.sortType = SortAsc
			cell.SetText(fmt.Sprintf("%s %s", t.columns[col], arrowDown))
		}
	} else {
		cell.SetText(t.columns[t.sortCol])
		newCell := t.GetCell(0, col)
		t.sortCol = col
		t.sortType = SortAsc
		newCell.SetText(fmt.Sprintf("%s %s", t.columns[col], arrowDown))
	}

	if t.sortFunc != nil {
		name := t.columns[t.sortCol]
		t.sortFunc(name, t.sortType)
	}
}
