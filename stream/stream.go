package stream

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/darshanime/netpeek/print"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
)

type HTTPStreamFactory struct{}

type httpReader struct {
	isClient bool
	data     []byte
	bytes    chan []byte
	stream   *httpStream
}

type httpStream struct {
	netFlow       gopacket.Flow
	transportFlow gopacket.Flow
	clientReader  httpReader
	serverReader  httpReader
	request       *http.Request
}

type AssemblerContext struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *AssemblerContext) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

func (h httpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	return true
}

func (h *httpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	fmt.Printf("Connection closed\n")
	return false
}

func (h *httpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, start, end, skip := sg.Info()
	fmt.Printf("info: dir %s, start %t, end %t, skip %t\n", dir, start, end, skip)
	length, _ := sg.Lengths()
	data := sg.Fetch(length)
	if dir == reassembly.TCPDirClientToServer {
		h.clientReader.bytes <- data
	} else {
		h.serverReader.bytes <- data
	}
	// fmt.Printf("got some data, (%d, %d): %s\n", length, saved, string(data))

}

func (h *httpReader) Read(p []byte) (int, error) {
	ok := true
	for ok && len(h.data) == 0 {
		h.data, ok = <-h.bytes
	}
	if !ok || len(h.data) == 0 {
		return 0, io.EOF
	}

	l := copy(p, h.data)
	h.data = h.data[l:]
	return l, nil
}

// New is required to statisfy the StreamFactory inferface
func (h *HTTPStreamFactory) New(netFlow, tcpFlow gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fmt.Printf("********* starting a new stream for: %s, %s\n", netFlow, tcpFlow)
	stream := &httpStream{
		netFlow:       netFlow,
		transportFlow: tcpFlow,
		clientReader:  httpReader{bytes: make(chan []byte), isClient: true},
		serverReader:  httpReader{bytes: make(chan []byte), isClient: false},
	}
	stream.clientReader.stream = stream
	stream.serverReader.stream = stream
	go stream.clientReader.read()
	go stream.serverReader.read()
	return stream
}

func (h *httpReader) read() {
	buf := bufio.NewReader(h)
	for {
		if h.isClient {
			req, err := http.ReadRequest(buf)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Printf("cannot read request, %s\n", err.Error())
			} else {
				h.stream.request = req
			}
		} else {
			resp, err := http.ReadResponse(buf, h.stream.request)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Printf("cannot read request, %s\n", err.Error())
			} else {
				print.Response(h.stream.request, resp)
			}
		}
	}
}
