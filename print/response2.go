package print

import (
	"fmt"
	"net/http"
	"os"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/stats"
)

func Response2(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
	fmt.Fprintf(os.Stdout, "\nRequest:\n%s", cui.RequestToString(req))
	fmt.Fprintf(os.Stdout, "\nResponse:\n%s", cui.ResponseToString(resp))
	fmt.Fprintf(os.Stdout, "\nPackets:\n%s", cui.PacketsToString(pktInfo))
}
