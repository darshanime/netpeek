package cui

import (
	"fmt"
	"net/http"

	"github.com/darshanime/netpeek/print"
	"github.com/darshanime/netpeek/stats"
	"github.com/jroimartin/gocui"
)

func PrintResponse(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo, name string) {
	maxX, maxY := g.Size()

	reqDetailName := "req->detail->" + name + "->req"
	respDetailName := "req->detail->" + name + "->resp"
	pktDetailName := "req->detail->" + name + "->pkt"

	// setting the request view
	logger.Printf("creating: " + "req->detail->" + name + "->req")
	v, err := g.SetView(reqDetailName, -1, 1, maxX/3, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}
		v.Autoscroll = false
		v.Frame = true
		v.Wrap = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(reqDetailName)
	fmt.Fprintln(v, print.RequestToString(req))

	// setting the response view
	logger.Printf("creating: " + respDetailName)
	v, err = g.SetView(respDetailName, maxX/3, 1, 2*maxX/3, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}
		v.Autoscroll = false
		v.Frame = true
		v.Wrap = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(respDetailName)
	fmt.Fprintln(v, print.ResponseToString(resp))

	// setting the packets view
	logger.Printf("creating: " + pktDetailName)
	v, err = g.SetView(pktDetailName, 2*maxX/3, 1, maxX, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}
		v.Autoscroll = false
		v.Frame = true
		v.Wrap = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(pktDetailName)
	fmt.Fprintln(v, print.PacketsToString(pktInfo))
}
