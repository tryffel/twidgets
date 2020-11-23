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
	"fmt"
	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
	"os"
	"tryffel.net/go/twidgets"
)

var app *cview.Application

func main() {
	app = cview.NewApplication()

	colors := twidgets.NavBarColors{
		Background:            tcell.Color235,
		BackgroundFocus:       tcell.Color235,
		ButtonBackground:      tcell.Color235,
		ButtonBackgroundFocus: tcell.Color23,
		Text:                  tcell.Color252,
		TextFocus:             tcell.Color253,
		Shortcut:              tcell.Color214,
		ShortcutFocus:         tcell.Color214,
	}
	navBar := twidgets.NewNavBar(&colors, done)

	btn := cview.NewButton("first")
	navBar.AddButton(btn, tcell.KeyF1)

	btn = cview.NewButton("second")
	navBar.AddButton(btn, tcell.KeyF2)

	btn = cview.NewButton("third")
	navBar.AddButton(btn, tcell.KeyF3)

	btn = cview.NewButton("fourth")
	navBar.AddButton(btn, tcell.KeyF4)

	app.SetRoot(navBar, true)
	app.Run()
}

func done(label string) {
	fmt.Printf("Got label %s", label)
	os.Exit(0)

}
