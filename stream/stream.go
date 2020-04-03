package stream

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/stats"
	"github.com/jroimartin/gocui"

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
	stats         streamStats
	cui           *gocui.Gui
}

type streamStats struct {
	startTime time.Time
	packets   []stats.PacketInfo
}

type AssemblerContext struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *AssemblerContext) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

func (h *httpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	captureInfo := ac.GetCaptureInfo()
	if *start {
		h.stats.startTime = captureInfo.Timestamp
	}
	pktInfo := stats.PacketInfo{
		FIN:           tcp.FIN,
		SYN:           tcp.SYN,
		RST:           tcp.RST,
		PSH:           tcp.PSH,
		ACK:           tcp.ACK,
		URG:           tcp.URG,
		ECE:           tcp.ECE,
		CWR:           tcp.CWR,
		NS:            tcp.NS,
		CaptureLength: captureInfo.CaptureLength,
		Timestamp:     captureInfo.Timestamp.Sub(h.stats.startTime),
		Dir:           dir.String(),
	}
	h.stats.packets = append(h.stats.packets, pktInfo)
	return true
}

func (h *httpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	// fmt.Println("Connection closed")
	return false
}

func (h *httpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, _ := sg.Info()
	length, _ := sg.Lengths()
	data := sg.Fetch(length)
	if dir == reassembly.TCPDirClientToServer {
		h.clientReader.bytes <- data
	} else {
		h.serverReader.bytes <- data
	}
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
	cui.AddConnection(netFlow, tcpFlow)
	stream := &httpStream{
		netFlow:       netFlow,
		transportFlow: tcpFlow,
		clientReader:  httpReader{bytes: make(chan []byte), isClient: true},
		serverReader:  httpReader{bytes: make(chan []byte), isClient: false},
	}
	stream.clientReader.stream = stream
	stream.serverReader.stream = stream
	go stream.clientReader.read(netFlow, tcpFlow)
	go stream.serverReader.read(netFlow, tcpFlow)
	return stream
}

func (h *httpReader) read(netFlow, tcpFlow gopacket.Flow) {
	buf := bufio.NewReader(h)
	for {
		if h.isClient {
			req, err := http.ReadRequest(buf)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Printf("~~cannot read request, %s\n", err.Error())
			} else {
				h.stream.request = req
			}
		} else {
			resp, err := http.ReadResponse(buf, h.stream.request)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Fprintln(os.Stderr, "cannot read request, %s\n", err.Error())
			} else {
				if *h.stream.useCui {
					cui.AddRequest(netFlow, tcpFlow, h.stream.request, resp, h.stream.stats.packets)
				} else {
					print.Response2(h.stream.request, resp, h.stream.stats.packets)
					h.stream.stats = streamStats{}
				}
			}
		}
	}
}
