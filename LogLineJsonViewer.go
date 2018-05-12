package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Jeffail/gabs"
	"github.com/jroimartin/gocui"
)

const logView string = "logView"
const detailView string = "detailsView"

type UIManager struct {
	notifier chan string
}

func main() {

	uiMgr := &UIManager{}
	uiMgr.notifier = make(chan string)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(uiMgr.layoutManager)

	g.SetKeybinding(logView, gocui.MouseWheelUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollUpwards(v)
	})

	g.SetKeybinding(logView, gocui.MouseWheelDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollDownwards(v)
	})

	g.SetKeybinding(detailView, gocui.MouseWheelUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollUpwards(v)
	})

	g.SetKeybinding(detailView, gocui.MouseWheelDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollDownwards(v)
	})

	g.SetKeybinding(logView, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {

		_, cy := v.Cursor()
		line, _ := v.Line(cy)

		if line == "" {
			return nil
		}

		detailsV, err := g.View(detailView)
		detailsV.Clear()
		container, err := gabs.ParseJSON([]byte(line))
		fmt.Fprint(detailsV, container.StringIndent("", "  "))

		return err
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func scrollUpwards(v *gocui.View) error {
	ox, oy := v.Origin()

	if oy == 0 {
		return nil
	}

	err := v.SetOrigin(ox, oy-1)
	ox, oy = v.Origin()

	return err
}

func scrollDownwards(v *gocui.View) error {
	ox, oy := v.Origin()
	err := v.SetOrigin(ox, oy+1)
	ox, oy = v.Origin()

	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (mgr *UIManager) layoutManager(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(logView, 0, 0, maxX-1, (maxY/2)+3); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log View"
		v.Wrap = true
	}

	if v, err := g.SetView(detailView, 0, (maxY/2)+4, maxX-1, (maxY - 2)); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Detail View"
		v.Wrap = true
	}

	mgr.logParse(g)
	return nil
}

func (mgr *UIManager) logParse(g *gocui.Gui) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) != 0 {

			logV, _ := g.View(logView)
			fmt.Fprintln(logV, text)
			fmt.Fprintln(logV, "")
		} else {
			break

		}

	}
}
