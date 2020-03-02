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
	"reflect"
	"testing"
)

type testItem struct {
	*tview.TextView
}

func (t *testItem) SetSelected(selection Selection) {}

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
		list       *ScrollList
		selected   int
		padding    int
		itemHeight int
		borders    bool
	}
	type args struct {
		x int
		y int
		w int
		h int
	}
	type out struct {
		rows        int
		visibleFrom int
		visibleTo   int
		wantGrid    []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		out    out
	}{
		{
			name: "one-row",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 2,
				borders:    true,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 4,
			},
			out: out{
				rows:        1,
				visibleFrom: 0,
				visibleTo:   0,
				wantGrid:    []int{2, -1},
			},
		},
		{
			name: "bordered",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 5,
				borders:    true,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 22,
			},
			out: out{
				rows:        3,
				visibleFrom: 0,
				visibleTo:   2,
				wantGrid:    []int{5, 1, 5, 1, 5, -1},
			},
		},
		{
			name: "bordered no bottom padding",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 5,
				borders:    true,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 23,
			},
			out: out{
				rows:        4,
				visibleFrom: 0,
				visibleTo:   3,
				wantGrid:    []int{5, 1, 5, 1, 5, 1, 5},
			},
		},
		{
			name: "bordered no bottom padding-3",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 5,
				borders:    true,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 25,
			},
			out: out{
				rows:        4,
				visibleFrom: 0,
				visibleTo:   3,
				wantGrid:    []int{5, 1, 5, 1, 5, 1, 5, -1},
			},
		},
		{
			name: "item-no-border",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 2,
				borders:    false,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 19,
			},
			out: out{
				rows:        6,
				visibleFrom: 0,
				visibleTo:   5,
				wantGrid:    []int{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, -1},
			},
		},
		{
			name: "2-height no padding",
			fields: fields{
				list:       NewScrollList(nil),
				selected:   0,
				padding:    1,
				itemHeight: 2,
				borders:    false,
			},
			args: args{
				x: 0,
				y: 0,
				w: 50,
				h: 10,
			},
			out: out{
				rows:        3,
				visibleFrom: 0,
				visibleTo:   2,
				wantGrid:    []int{2, 1, 2, 1, 2, -1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.fields.list.selected = tt.fields.selected
			tt.fields.list.Padding = tt.fields.padding
			tt.fields.list.ItemHeight = tt.fields.itemHeight

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

			if !reflect.DeepEqual(tt.fields.list.gridRows, tt.out.wantGrid) {
				t.Errorf("scoll_list.updateGrid() grid rows: got %v, expected %v",
					tt.fields.list.gridRows, tt.out.wantGrid)
			}
		})
	}
}
