package cui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

var g *gocui.Gui

var connMap map[gopacket.Flow]int

func InitCui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	connMap = map[gopacket.Flow]int{}

	if err != nil {
		log.Panicln(err)
	}
	g = gui

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, actionViewConnectionsUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, actionViewConnectionsDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, actionEnterKey); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("title", -1, -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false
		v.BgColor = gocui.ColorDefault | gocui.AttrReverse
		v.FgColor = gocui.ColorDefault | gocui.AttrReverse

		// Content
		fmt.Fprintln(v, "⣿ NetPeek")
	}

	if v, err := g.SetView("conns", -1, 1, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 2)

		viewConnsAddLine(v, "SRC", "DST")
		fmt.Fprintln(v, strings.Repeat("─", maxX))

		g.SetCurrentView(v.Name())
		// go viewConnectionsWithAutoRefresh(g)
	}

	if v, err := g.SetView("status", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite

		// Content
		changeStatusContext(g, "C")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	os.Exit(0)
	return nil
}

func viewConnsAddLine(v *gocui.View, src, dst string) {
	line := pad.Right(src, 30, " ") + pad.Right(dst, 30, " ")
	fmt.Fprintln(v, line)
}

func changeStatusContext(g *gocui.Gui, c string) error {
	lMaxX, _ := g.Size()
	v, err := g.View("status")
	if err != nil {
		return err
	}

	v.Clear()

	i := lMaxX + 4
	b := ""

	switch c {
	case "C":
		i = 150 + i
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("Enter") + " View Requests   "
	case "SE":
		i = i + 100
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("Enter") + " Select   "
	case "SL":
		i = i + 100
		b = b + frameText("q") + " Back   "
		b = b + frameText("Enter") + " Select   "
	}
	b = b + frameText("CTRL+C") + " Exit"

	fmt.Fprintln(v, pad.Right(b, i, " "))

	return nil
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
