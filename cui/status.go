package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

func updateStatus(g *gocui.Gui, c string) error {
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
	case "RL":
		i = i + 100
		b = b + frameText("←") + " Back   "
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("Enter") + " View Details   "
	case "RD":
		i = i + 100
		b = b + frameText("←") + " Back   "
	}
	b = b + frameText("CTRL+C") + " Exit"
	fmt.Fprintln(v, pad.Right(b, i, " "))
	return nil
}
