package stream

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/print"
	"github.com/darshanime/netpeek/stats"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
)

type Protocol int

type HTTPStreamFactory struct {
	UseCui   *bool
	Protocol *string
	Logger   *log.Logger
}

const (
	HTTP Protocol = iota
	Dump
	Drain
)

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
	protocol      Protocol
	logger        *log.Logger
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
	h.logger.Printf("closing old connection, %s", connDir(h.netFlow, h.transportFlow))
	if h.transportFlow.Src().String() == "443" || h.transportFlow.Dst().String() == "443" {
		if *h.useCui {
			cui.AddRequest(h.netFlow, h.transportFlow, nil, nil, h.stats.packets)
		} else {
			print.Response(nil, nil, h.stats.packets)
		}
		h.clientReader.stream.stats = streamStats{}
		h.serverReader.stream.stats = streamStats{}
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
	h.Logger.Printf("\nadding new connection, %s", connDir(netFlow, tcpFlow))
	stream := &httpStream{
		netFlow:       netFlow,
		transportFlow: tcpFlow,
		clientReader:  httpReader{bytes: make(chan []byte), isClient: true},
		serverReader:  httpReader{bytes: make(chan []byte), isClient: false},
		useCui:        h.UseCui,
		logger:        h.Logger,
	}
	stream.clientReader.stream = stream
	stream.serverReader.stream = stream
	stream.protocol = getProtocol(*h.Protocol)
	go stream.clientReader.read()
	go stream.serverReader.read()
	return stream
}

func (h *httpReader) read() {
	switch h.stream.protocol {
	case HTTP:
		if h.isClient {
			h.stream.logger.Printf("starting http request reader")
			go readHTTPRequest(h)
		} else {
			h.stream.logger.Printf("starting http response reader")
			go readHTTPResponse(h)
		}
	case Drain:
		go drainPackets(h)
	case Dump:
		go dumpPackets(h)
	}
}

func readHTTPResponse(h *httpReader) {
	buf := bufio.NewReader(h)
	for {
		resp, err := http.ReadResponse(buf, h.stream.request)
		h.stream.logger.Printf("read response %v", err)
		if err == io.EOF {
			h.stream.logger.Printf("stopped reading response, got EOF, %s", err.Error())
			return
		} else if err != nil {
			h.stream.logger.Printf("cannot read response, %s", err.Error())
		} else {
			if *h.stream.useCui {
				cui.AddRequest(h.stream.netFlow, h.stream.transportFlow, h.stream.request, resp, h.stream.stats.packets)
			} else {
				print.Response(h.stream.request, resp, h.stream.stats.packets)
				h.stream.stats = streamStats{}
			}
		}
	}
}

func readHTTPRequest(h *httpReader) {
	buf := bufio.NewReader(h)
	for {
		req, err := http.ReadRequest(buf)
		h.stream.logger.Printf("read request %v", err)
		if err == io.EOF {
			h.stream.logger.Printf("stopped reading request, got EOF, %s", err.Error())
			return
		} else if err != nil {
			h.stream.logger.Printf("cannot read request, %s", err.Error())
		} else {
			h.stream.request = req
		}
	}
}

func drainPackets(h *httpReader) {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			if len(h.stream.stats.packets) != 0 {
				if *h.stream.useCui {
					cui.AddRequest(h.stream.netFlow, h.stream.transportFlow, nil, nil, h.stream.stats.packets)
				} else {
					print.Response(nil, nil, h.stream.stats.packets)
				}
				h.stream.stats = streamStats{}
			}
		}
	}
}

func dumpPackets(h *httpReader) {
	buf := bufio.NewReader(h)
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			if len(h.stream.stats.packets) != 0 {
				io.CopyN(h.stream.logger.Writer(), buf, int64(buf.Size()))
			}
			h.stream.stats = streamStats{}
		}
	}
}

func connDir(netflow, tcpflow gopacket.Flow) string {
	return netflow.Src().String() + ":" + tcpflow.Src().String() + "-->" + netflow.Dst().String() + ":" + tcpflow.Dst().String()
}

func getProtocol(protocol string) Protocol {
	switch protocol {
	case "http":
		return HTTP
	case "drain":
	case "https":
		return Drain
	case "dump":
		return Dump
	}
	panic(fmt.Sprintf("unknown protocol: %s", protocol))
}
