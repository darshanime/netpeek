package print

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/darshanime/netpeek/stats"
)

func Response(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
	var reqStr, respStr, pktStr string
	if req != nil {
		reqStr = RequestToString(req)
	} else {
		reqStr = "**** no req\n"
	}
	if resp != nil {
		respStr = ResponseToString(resp)
	} else {
		respStr = "**** no resp\n"
	}
	pktStr = PacketsToString(pktInfo)
	fmt.Fprintf(os.Stdout, "\nRequest:\n%s\nResponse:\n%s\n\nPackets:\n%s", reqStr, respStr, pktStr)
}

func RequestToString(req *http.Request) string {
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

func ResponseToString(resp *http.Response) string {
	var str strings.Builder
	str.WriteString("Response code: " + resp.Status + "\n")
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

func PacketsToString(pktInfo []stats.PacketInfo) string {
	var str strings.Builder
	for _, pkt := range pktInfo {
		str.WriteString(pkt.String() + "\n")
	}
	return str.String()
}
