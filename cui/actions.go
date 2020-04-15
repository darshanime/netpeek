package cui

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

func actionGlobalQuit(g *gocui.Gui, v *gocui.View) error {
	os.Exit(0)
	return nil
}

func actionEnterKey(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "conns" {
		conn, err := getSelectedConnection(g)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "setting on top: "+"reqs->list"+conn)
		updateStatus(g, "RL")
		g.SetCurrentView("reqs->list" + conn)
		g.SetViewOnTop("reqs->list" + conn)
		return nil
	} else if strings.Contains(v.Name(), "reqs->list") {
		conn, err := getSelectedRequest()
		if err != nil {
			return err
		}
		updateStatus(g, "RD")
		g.SetCurrentView(conn + "->pkt")
		g.SetViewOnTop(conn + "->pkt")
		g.SetCurrentView(conn + "->resp")
		g.SetViewOnTop(conn + "->resp")
		g.SetCurrentView(conn + "->req")
		g.SetViewOnTop(conn + "->req")
		fmt.Fprintln(os.Stderr, "setting on top: "+conn+"->pkt")
		return nil
	}
	return nil
}

func actionArrowLeftKey(g *gocui.Gui, v *gocui.View) error {
	if strings.Contains(v.Name(), "reqs->list") {
		updateStatus(g, "C")
		g.SetCurrentView("conns")
		g.SetViewOnTop("conns")
		return nil
	} else if strings.Contains(v.Name(), "req->detail") {
		splitReqList := strings.Split(v.Name(), "->")
		recListViewName := "reqs->list->" + splitReqList[3] + "->" + splitReqList[4]
		fmt.Fprintln(os.Stderr, "trying to go back: "+recListViewName)
		g.SetCurrentView(recListViewName)
		g.SetViewOnTop(recListViewName)
		return nil
	}
	return nil
}

func actionGlobalArrowDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
	return nil
}

func actionGlobalArrowUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 2)
	return nil
}

func getSelectedRequest() (string, error) {
	v := g.CurrentView()
	l, err := getViewLine(g, v)
	if err != nil {
		return "", err
	}
	number := getRequestNumberFromLine(l)
	splitReqList := strings.Split(v.Name(), "->")
	return "req->detail->" + number + "->" + splitReqList[2] + "->" + splitReqList[3], nil
}

func getSelectedConnection(g *gocui.Gui) (string, error) {
	v, err := g.View("conns")
	if err != nil {
		return "", err
	}
	l, err := getViewLine(g, v)
	if err != nil {
		return "", err
	}
	p := getConnectionNameFromLine(l)
	return p, nil
}

func getViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, err
}

func getConnectionNameFromLine(line string) string {
	splitLine := strings.Split(line, " ")
	result := ""
	for _, substr := range splitLine {
		if substr != "" {
			result += "->" + substr
		}
	}
	return result
}

func getRequestNumberFromLine(line string) string {
	splitLine := strings.Split(line, " ")
	return splitLine[0]
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
