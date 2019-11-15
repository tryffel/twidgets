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
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

//Modal interface creates a modal that overlaps other views and get's destroyed when it's ready

//Modal interface is primitive that is drawn on top of other views and get's destroyed when it's ready
type Modal interface {
	//Primitive
	tview.Primitive
	//SetDoneFunc sets function that get's called when modal wants to close itself
	SetDoneFunc(doneFunc func())
	//Setvisible tells modal to show or hide itself
	SetVisible(visible bool)
}

/*ModalLayout is a grid layout that draws modal on center of grid.
To add ordinary items to layout, get grid with Grid() function. Layout consists of 6 columns and rows
defined as (2, -1, -2, -2, -1, 2). Modal is drawn on middle 4 cells

Use AddModal and RemoveModal to manage modals. Only single modal can be shown at a time.
 */
type ModalLayout struct {
	grid       *tview.Grid
	hasModal   bool
	customGrid bool
	modal      Modal

	//Default grid col/row weights
	gridAxis   []int
}

//NewModalLayout creates new modal layout and returns it
func NewModalLayout() *ModalLayout {
	m := &ModalLayout{
		grid:       tview.NewGrid(),
		hasModal:   false,
		customGrid: false,
		modal:      nil,
		gridAxis:   nil,
	}

	/*
	Put modal to rows/cols 3-4
	Changing these requires also changing AddColumn()-> grid.AddItem indices.
	 */
	m.gridAxis = []int{2, -1, -2, -2, -1, 2}

	m.grid.SetRows(m.gridAxis...)
	m.grid.SetColumns(m.gridAxis...)
	m.grid.SetMinSize(2, 2)

	return m
}

func (m *ModalLayout) Draw(screen tcell.Screen) {
	m.grid.Draw(screen)
}

func (m *ModalLayout) GetRect() (int, int, int, int) {
	return m.grid.GetRect()
}

func (m *ModalLayout) SetRect(x, y, width, height int) {
	m.grid.SetRect(x, y, width, height)
}

func (m *ModalLayout) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.grid.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

	})
}

func (m *ModalLayout) Focus(delegate func(p tview.Primitive)) {
	m.grid.Focus(delegate)
}

func (m *ModalLayout) Blur() {
	m.grid.Blur()
}

func (m *ModalLayout) GetFocusable() tview.Focusable {
	return m.grid.GetFocusable()
}

func (m *ModalLayout) Grid() *tview.Grid {
	return m.grid
}

//AddModal adds modal to center of layout
// lockSize flag defines whether modal size should be locked or dynamic.
func (m *ModalLayout) AddModal(modal Modal, height, width uint, lockSize bool) {
	if m.hasModal {
		return
	}
	if !lockSize {
		m.customGrid = false
		m.grid.AddItem(modal, 2, 2, 2, 2, 8, 30, true)
	} else {
		m.customGrid = true
		x := make([]int, len(m.gridAxis))
		y := make([]int, len(m.gridAxis))
		copy(x, m.gridAxis)
		copy(y, m.gridAxis)
		x[2] = int(width / 2)
		x[3] = x[2]
		y[2] = int(height / 2)
		y[3] = y[2]
		m.grid.SetRows(y...)
		m.grid.SetColumns(x...)
		m.grid.AddItem(modal, 2, 2, 2, 2, int(height), int(width), true)
	}
	m.hasModal = true
	m.modal = modal
}

//RemoveModal removes modal
func (m *ModalLayout) RemoveModal(modal Modal) {
	if m.hasModal {
		modal.SetVisible(false)
		m.grid.RemoveItem(modal)
		m.hasModal = false
		m.modal = nil
		if m.customGrid {
			m.grid.SetRows(m.gridAxis...)
			m.grid.SetColumns(m.gridAxis...)
			m.customGrid = false
		}
	}
}
