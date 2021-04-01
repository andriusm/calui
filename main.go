package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	week Week
)

const dayHiColor = tcell.Color118
const eventHiColor = tcell.Color87

type Week struct {
	days    []*Day
	current int
}

type Day struct {
	title   string
	box     *tview.Flex
	events  []Event
	current int
}

type Event struct {
	title string
	box   *tview.TextView
}

func (w *Week) AddDay(title string) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle("[ " + title + " ]")
	w.days = append(w.days, &Day{title: title, box: flex})
}

func (d *Day) AddEvent(title string) {
	text := tview.NewTextView()
	text.SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(title)
	text.SetScrollable(true)
	//text.SetTitle(" event ")

	d.events = append(d.events, Event{title: title, box: text})
	d.box.AddItem(text, 5, 1, false)
}

func (d *Day) NextEvent() {
	if d.current+1 < len(d.events) {
		d.events[d.current].box.SetBorderColor(tcell.ColorWhite)
		d.current += 1
		d.events[d.current].box.SetBorderColor(eventHiColor)
	}
}

func (d *Day) PrevEvent() {
	if d.current-1 >= 0 {
		d.events[d.current].box.SetBorderColor(tcell.ColorWhite)
		d.current -= 1
		d.events[d.current].box.SetBorderColor(eventHiColor)
	}
}

func (w *Week) PrevDay() {
	if w.current-1 >= 0 {
		d := w.days[w.current]
		d.box.SetBorderColor(tcell.ColorWhite)
		d.events[d.current].box.SetBorderColor(tcell.ColorWhite)

		w.current -= 1

		d = w.days[w.current]
		d.box.SetBorderColor(dayHiColor)
		d.events[d.current].box.SetBorderColor(eventHiColor)
	}
}

func (w *Week) NextDay() {
	if w.current+1 < len(week.days) {
		d := w.days[w.current]
		d.box.SetBorderColor(tcell.ColorWhite)
		d.events[d.current].box.SetBorderColor(tcell.ColorWhite)

		w.current += 1

		d = w.days[w.current]
		d.box.SetBorderColor(dayHiColor)
		d.events[d.current].box.SetBorderColor(eventHiColor)
	}
}

func myEvents(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'q' {
		app.Stop()
	}
	if event.Rune() == 'j' {
		day := week.days[week.current]
		day.NextEvent()
	}
	if event.Rune() == 'k' {
		day := week.days[week.current]
		day.PrevEvent()
	}
	if event.Rune() == 'h' {
		week.PrevDay()
	}
	if event.Rune() == 'l' {
		week.NextDay()
	}
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
	}

	flex := tview.NewFlex()
	for _, d := range week.days {
		flex.AddItem(d.box, 0, 1, false)
	}

	app = tview.NewApplication().SetRoot(flex, true).EnableMouse(true)
	app.SetInputCapture(myEvents)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
