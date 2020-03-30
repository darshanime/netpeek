package cui

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

func actionEnterKey(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "conns" {
		conn, err := getSelectedConnection(g)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "setting on top: "+"reqs->list"+conn)
		g.SetCurrentView("reqs->list" + conn)
		g.SetViewOnTop("reqs->list" + conn)
		return nil
	} else if strings.Contains(v.Name(), "reqs->list") {
		conn, err := getSelectedRequest()
		if err != nil {
			return err
		}
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
