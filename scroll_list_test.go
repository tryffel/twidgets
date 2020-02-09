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
	"fmt"
	"github.com/rivo/tview"
	"testing"
)

type testItem struct {
	*tview.TextView
}

func (t *testItem) SetSelected(bool) {}


func TestScrollList_updateGrid(t *testing.T) {
	// Test grid size & num of elements is updated correctly

	createItems := 10
	items := make([]*testItem, createItems)

	for i := 0; i < createItems; i++ {
		item := &testItem{tview.NewTextView()}
		item.SetBorder(false)

		item.SetText(fmt.Sprintf("Item %d\n2nd row", i))
		items[i] = item
	}

	type fields struct {
		list *ScrollList
		selected int
		padding int
		itemHeight int
		borders bool
	}
	type args struct {
		x int
		y int
		w int
		h int
	}
	type out struct {
		rows int
		visibleFrom int
		visibleTo int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		out    out
	}{
		{
			name: "item-border",
			fields:fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 5,
				borders:    true,
			},
			args:args{
				x: 0,
				y: 0,
				w: 50,
				h: 22,
			},
			out:out{
				rows:        3,
				visibleFrom: 0,
				visibleTo:   2,
			},

		},
		{
			name: "item-no-border",
			fields:fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 2,
				borders:    false,
			},
			args:args{
				x: 0,
				y: 0,
				w: 50,
				h: 19,
			},
			out:out{
				rows:        5,
				visibleFrom: 0,
				visibleTo:   4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {


			tt.fields.list.selected = tt.fields.selected
			tt.fields.list.Padding= tt.fields.padding
			tt.fields.list.ItemHeight= tt.fields.itemHeight

			for i := 0; i < createItems; i++ {
				items[i].SetBorder(tt.fields.borders)
				tt.fields.list.AddItem(items[i])
			}

			args := tt.args
			tt.fields.list.updateGrid(args.x, args.y, args.w, args.h)

			if tt.fields.list.rows != tt.out.rows {
				t.Errorf("scroll_list.updateGrid() rows: got %d, expected %d", tt.fields.list.rows, tt.out.rows)
			}

			if tt.fields.list.visibleFrom != tt.out.visibleFrom {
				t.Errorf("scroll_list.updateGrid() visible from: got %d, expected %d",
					tt.fields.list.visibleFrom, tt.out.visibleFrom)
			}

			if tt.fields.list.visibleTo != tt.out.visibleTo {
				t.Errorf("scroll_list.updateGrid() visible to: got %d, expected %d",
					tt.fields.list.visibleTo, tt.out.visibleTo)
			}

		})
	}
}