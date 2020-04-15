package cui

import (
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

var g *gocui.Gui

var connMap map[gopacket.Flow]int

type Key struct {
	viewname string
	key      interface{}
	handler  func(*gocui.Gui, *gocui.View) error
}

var keys []Key = []Key{
	Key{"", gocui.KeyCtrlC, actionGlobalQuit},
	Key{"", gocui.KeyArrowUp, actionGlobalArrowUp},
	Key{"", gocui.KeyArrowDown, actionGlobalArrowDown},
	Key{"", gocui.KeyEnter, actionEnterKey},
	Key{"", gocui.KeyArrowRight, actionEnterKey},
	Key{"", gocui.KeyArrowLeft, actionArrowLeftKey},
}

func InitCui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	connMap = map[gopacket.Flow]int{}

	if err != nil {
		panic(fmt.Sprintf("could not init new gui: %s\n", err.Error()))
	}

	g = gui

	g.SetManagerFunc(layout)

	for _, key := range keys {
		if err := g.SetKeybinding(key.viewname, key.key, gocui.ModNone, key.handler); err != nil {
			panic(fmt.Sprintf("could not set key bindings: %s\n", err.Error()))
		}
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(fmt.Sprintf("received error in main loop: %s\n", err.Error()))
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("title", -1, -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.BgColor = gocui.ColorDefault | gocui.AttrReverse
		v.FgColor = gocui.ColorDefault | gocui.AttrReverse

		fmt.Fprintln(v, "⣿ NetPeek")
	}

	if v, err := g.SetView("conns", -1, 1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 2)

		viewConnsAddLine(v, "SRC", "DST")
		fmt.Fprintln(v, strings.Repeat("─", maxX))

		g.SetCurrentView(v.Name())
	}

	if v, err := g.SetView("status", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite

		updateStatus(g, "C")
	}
	return nil
}

func viewConnsAddLine(v *gocui.View, src, dst string) {
	line := pad.Right(src, 30, " ") + pad.Right(dst, 30, " ")
	fmt.Fprintln(v, line)
}

// StringFormatBoth fg and bg colors
// Thanks https://github.com/mephux/komanda-cli/blob/master/komanda/color/color.go
func stringFormatBoth(fg, bg int, str string, args []string) string {
	return fmt.Sprintf("\x1b[48;5;%dm\x1b[38;5;%d;%sm%s\x1b[0m", bg, fg, strings.Join(args, ";"), str)
}

// Frame text with colors
func frameText(text string) string {
	return stringFormatBoth(15, 0, text, []string{"1"})
}

func moveViewCursorUp(g *gocui.Gui, v *gocui.View, dY int) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if cy > dY {
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func moveViewCursorDown(g *gocui.Gui, v *gocui.View, allowEmpty bool) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, err := getNextViewLine(g, v)
	if err != nil {
		return err
	}
	if !allowEmpty && nextLine == "" {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func getNextViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}
