package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	week Week
)

const (
	hiColor   = tcell.ColorGreenYellow
	loColor   = tcell.ColorDimGray
	backColor = tcell.ColorGray
)

type Week struct {
	days       []*Day
	currentDay int
}

type Day struct {
	title        string
	box          *tview.Flex
	events       []Event
	currentEvent int
}

type Event struct {
	title string
	box   *tview.TextView
}

func (w *Week) AddDay(title string) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle("[ " + title + " ]")
	flex.SetBackgroundColor(backColor)
	w.days = append(w.days, &Day{title: title, box: flex})
}

func (d *Day) AddEvent(title string) {
	text := tview.NewTextView()
	text.SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(title)
	text.SetScrollable(true)
	text.SetBackgroundColor(backColor)

	d.events = append(d.events, Event{title: title, box: text})
	d.box.AddItem(text, 5, 1, false)
}

func (d *Day) NextEvent() {
	if d.currentEvent+1 < len(d.events) {
		d.currentEvent += 1
	}
}

func (d *Day) PrevEvent() {
	if d.currentEvent-1 >= 0 {
		d.currentEvent -= 1
	}
}

func (w *Week) PrevDay() {
	if w.currentDay-1 >= 0 {
		w.currentDay -= 1
	}
}

func (w *Week) NextDay() {
	if w.currentDay+1 < len(week.days) {
		w.currentDay += 1
	}
}

func (w *Week) CurrentDay() *Day {
	return w.days[w.currentDay]
}

func (w *Week) Redraw() {
	for i, d := range w.days {
		if i == w.currentDay {
			d.box.SetBorderColor(hiColor)
		} else {
			d.box.SetBorderColor(loColor)
		}

		for j, e := range d.events {
			if j == d.currentEvent && i == w.currentDay {
				e.box.SetBorderColor(hiColor)
			} else {
				e.box.SetBorderColor(loColor)
			}
		}
	}
}

func keyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		app.Stop()
	case 'j':
		week.CurrentDay().NextEvent()
	case 'k':
		week.CurrentDay().PrevEvent()
	case 'h':
		week.PrevDay()
	case 'l':
		week.NextDay()
	}

	week.Redraw()

	return event
}

func init() {
	week = Week{}
	week.AddDay("Monday")
	week.AddDay("Tuesday")
	week.AddDay("Wednesday")
	week.AddDay("Thursday")
	week.AddDay("Friday")
}

func main() {
	events := getEvents()
	for i, d := range events {
		for _, e := range d {
			week.days[i].AddEvent(e)
		}

		filler := tview.NewBox().SetBorder(false).SetBackgroundColor(backColor)
		week.days[i].box.AddItem(filler, 0, 1, false)
	}

	flex := tview.NewFlex()
	for _, d := range week.days {
		flex.AddItem(d.box, 0, 1, false)
	}

	app = tview.NewApplication().SetRoot(flex, true).EnableMouse(true)
	app.SetInputCapture(keyEvents)

	week.Redraw()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
