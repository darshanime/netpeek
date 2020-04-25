package print

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/darshanime/netpeek/stats"
)

func Output(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
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

func ResponseToString(resp *http.Response) string {
	var str strings.Builder
	str.WriteString("Response code: " + resp.Status + "\n")
	for key, val := range resp.Header {
		str.WriteString(key + ": " + strings.Join(val, ",") + "\n")
	}
	str.WriteString("\n")

	var reader io.ReadCloser

	if contentEncoding := resp.Header.Get("Content-Encoding"); contentEncoding == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			reader = resp.Body
			defer gzipReader.Close()
		} else {
			reader = gzipReader
		}
	} else {
		reader = resp.Body
	}
	defer resp.Body.Close()

	rcBody, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
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
