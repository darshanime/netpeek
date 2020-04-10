package print

import (
	"fmt"
	"net/http"
	"os"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/stats"
)

func Response2(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
	var reqStr, respStr, pktStr string
	if req != nil {
		reqStr = cui.RequestToString(req)
	} else {
		reqStr = "**** no req\n"
	}
	if resp != nil {
		respStr = cui.ResponseToString(resp)
	} else {
		respStr = "**** no resp\n"
	}
	pktStr = cui.PacketsToString(pktInfo)
	fmt.Fprintf(os.Stdout, "\nRequest:\n%s\nResponse:\n%s\n\nPackets:\n%s", reqStr, respStr, pktStr)
}
