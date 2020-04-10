package cui

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/darshanime/netpeek/stats"
	"github.com/google/gopacket"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

var requestCounter map[gopacket.Flow]map[gopacket.Flow]int

func init() {
	requestCounter = make(map[gopacket.Flow]map[gopacket.Flow]int)
}

func AddRequest(netflow, tcpflow gopacket.Flow, req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
	conn := getRequestName(netflow, tcpflow)
	if req == nil {
		fmt.Fprintln(os.Stderr, "not creating: reqs->list"+conn)
		return
	}
	maxX, maxY := g.Size()

	fmt.Fprintln(os.Stderr, "creating: reqs->list"+conn)

	v, err := g.SetView("reqs->list"+conn, -1, 1, maxX, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		panic("error with view")
	}
	if err == gocui.ErrUnknownView {
		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 2)
		requestsListAddLine(v, "", "Request", "Code", "#pkts", "Latency")
		fmt.Fprintln(v, strings.Repeat("─", maxX))
	}

	g.SetViewOnBottom("reqs->list" + conn)

	if _, ok := requestCounter[netflow]; !ok {
		requestCounter[netflow] = map[gopacket.Flow]int{tcpflow: 0}
	}
	requestCounter[netflow][tcpflow]++
	requestNum := strconv.Itoa(requestCounter[netflow][tcpflow])
	numPkts := len(pktInfo)

	requestsListAddLine(v, requestNum, fmt.Sprintf("%s %s", req.Method, req.URL.String()), strconv.Itoa(resp.StatusCode), strconv.Itoa(numPkts), pktInfo[numPkts-1].Timestamp.String())
	PrintResponse(req, resp, pktInfo, requestNum+conn)
}

func getRequestName(netflow, tcpflow gopacket.Flow) string {
	return "->" + netflow.Src().String() + ":" + tcpflow.Src().String() + "->" + netflow.Dst().String() + ":" + tcpflow.Dst().String()
}

func requestsListAddLine(v *gocui.View, sno, request, code, numPkts, latency string) {
	line := pad.Right(sno, 10, " ") + pad.Right(request, 20, " ") + pad.Right(code, 20, " ") + pad.Right(numPkts, 20, " ") + pad.Right(latency, 20, " ")
	fmt.Fprintln(v, line)
}
