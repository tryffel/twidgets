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
	"tryffel.net/go/twidgets"
)

func main() {
	app := tview.NewApplication()

	table := twidgets.NewTable()
	app.SetRoot(table, true)

	table.SetAddCellFunc(addCell)
	table.SetShowIndex(true)
	table.SetColumns([]string{"Name", "Description"})
	table.SetSort(1, twidgets.SortAsc)

	table.AddRow(0, "A", "first sample")
	table.AddRow(1, "B", "seconds sample")
	table.AddRow(2, "C", "third sample")
	table.AddRow(3, "D", "fourth sample")

	sortFunc := func(col string, sort twidgets.Sort) {
		if col == "Name" && sort == twidgets.SortAsc {
			table.Clear(false)
			table.AddRow(0, "A", "first sample")
			table.AddRow(1, "B", "second sample")
			table.AddRow(2, "C", "third sample")
			table.AddRow(3, "D", "fourth sample")
		} else if col == "Name" && sort == twidgets.SortDesc {
			table.Clear(false)
			table.AddRow(3, "A", "first sample")
			table.AddRow(2, "B", "second sample")
			table.AddRow(1, "C", "third sample")
			table.AddRow(0, "D", "fourth sample")
		} else if col == "Description" && sort == twidgets.SortAsc {
			table.Clear(false)
			table.AddRow(0, "A", "first sample")
			table.AddRow(2, "B", "second sample")
			table.AddRow(3, "C", "third sample")
			table.AddRow(1, "D", "fourth sample")
		} else if col == "Description" && sort == twidgets.SortDesc {
		table.Clear(false)
		table.AddRow(3, "A", "first sample")
		table.AddRow(1, "B", "second sample")
		table.AddRow(0, "C", "third sample")
		table.AddRow(2, "D", "fourth sample")
	}
	}
	table.SetSortFunc(sortFunc)

	app.Run()
}


func addCell(cell *tview.TableCell, header bool, col int) {
	if header {
		cell.SetBackgroundColor(tcell.ColorGreen)
	} else {
		// Change color for every 2nd row
		if col % 2 == 0 {
			cell.SetBackgroundColor(tcell.ColorGray)
		} else {
			cell.SetBackgroundColor(tcell.ColorDarkGray)
		}
	}
}