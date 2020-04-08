package stream

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/print"
	"github.com/darshanime/netpeek/stats"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
)

type HTTPStreamFactory struct {
	UseCui *bool
}

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
	useCui        *bool
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
	if *start || h.stats.startTime.IsZero() {
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
	fmt.Fprintln(os.Stderr, fmt.Sprintf("\nclosing old connection, %s\n", connDir(h.netFlow, h.transportFlow)))
	if h.transportFlow.Src().String() == "443" || h.transportFlow.Dst().String() == "443" {
		print.Response2(nil, nil, h.stats.packets)
	}
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
	if *h.UseCui {
		cui.AddConnection(netFlow, tcpFlow)
	} else {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("\nadding new connection, %s\n", connDir(netFlow, tcpFlow)))
	}
	stream := &httpStream{
		netFlow:       netFlow,
		transportFlow: tcpFlow,
		clientReader:  httpReader{bytes: make(chan []byte), isClient: true},
		serverReader:  httpReader{bytes: make(chan []byte), isClient: false},
		useCui:        h.UseCui,
	}
	stream.clientReader.stream = stream
	stream.serverReader.stream = stream
	go stream.clientReader.read()
	go stream.serverReader.read()
	return stream
}

func (h *httpReader) read() {
	buf := bufio.NewReader(h)
	if isSSL(h.stream.transportFlow) {
		go drainPackets(h)
		io.Copy(ioutil.Discard, buf)
	}
	for {
		if h.isClient {
			req, err := http.ReadRequest(buf)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Fprintln(os.Stderr, "cannot read request, %s\n", err.Error())
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
					cui.AddRequest(h.stream.netFlow, h.stream.transportFlow, h.stream.request, resp, h.stream.stats.packets)
				} else {
					print.Response2(h.stream.request, resp, h.stream.stats.packets)
					h.stream.stats = streamStats{}
				}
			}
		}
	}
}

func isSSL(tcpflow gopacket.Flow) bool {
	return tcpflow.Src().String() == "443" || tcpflow.Dst().String() == "443"
}

func drainPackets(h *httpReader) {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			if len(h.stream.stats.packets) != 0 {
				print.Response2(nil, nil, h.stream.stats.packets)
			}
			h.stream.stats = streamStats{}
		}
	}
}

func connDir(netflow, tcpflow gopacket.Flow) string {
	return netflow.Src().String() + ":" + tcpflow.Src().String() + "-->" + netflow.Dst().String() + ":" + tcpflow.Dst().String()
}
