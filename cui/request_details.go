package cui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/darshanime/netpeek/stats"
	"github.com/jroimartin/gocui"
)

func PrintResponse(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo, name string) {
	maxX, maxY := g.Size()
	fmt.Fprintln(os.Stderr, "creating: "+"req->detail->"+name+"->req")

	reqDetailName := "req->detail->" + name + "->req"
	respDetailName := "req->detail->" + name + "->resp"
	pktDetailName := "req->detail->" + name + "->pkt"

	// setting the request view
	v, err := g.SetView(reqDetailName, -1, 1, maxX/3, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}

		v.Frame = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(reqDetailName)
	fmt.Fprintln(v, requestToString(req))

	// setting the response view
	fmt.Fprintln(os.Stderr, "creating: "+respDetailName)
	v, err = g.SetView(respDetailName, maxX/3, 1, 2*maxX/3, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}

		v.Frame = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(respDetailName)
	fmt.Fprintln(v, responseToString(resp))

	// setting the packets view
	fmt.Fprintln(os.Stderr, "creating: "+pktDetailName)
	v, err = g.SetView(pktDetailName, 2*maxX/3, 1, maxX, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			panic("panic in set view")
		}
		v.Frame = true
	}
	v.SetCursor(0, 2)
	g.SetViewOnBottom(pktDetailName)
	fmt.Fprintln(v, packetsToString(pktInfo))
}

func getRequestDetailView(req *http.Request, suffix string) string {
	return req.URL.String() + suffix
}

func requestToString(req *http.Request) string {
	var str strings.Builder
	str.WriteString(req.Method + " " + req.URL.String() + "\n")
	for key, val := range req.Header {
		str.WriteString(key + ": " + strings.Join(val, ",") + "\n")
	}
	str.WriteString("\n")
	rcBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
	defer req.Body.Close()
	str.Write(rcBody)
	return str.String()
}

func responseToString(resp *http.Response) string {
	var str strings.Builder
	for key, val := range resp.Header {
		str.WriteString(key + ": " + strings.Join(val, ",") + "\n")
	}
	str.WriteString("\n")
	rcBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
	defer resp.Body.Close()
	str.Write(rcBody)
	return str.String()
}

func packetsToString(pktInfo []stats.PacketInfo) string {
	var str strings.Builder
	for _, pkt := range pktInfo {
		str.WriteString(pkt.String() + "\n")
	}
	return str.String()
}
