package stream

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
)

type HTTPStreamFactory struct{}

type httpReader struct {
	bytes chan []byte
}

type httpStream struct {
	netFlow       gopacket.Flow
	transportFlow gopacket.Flow
	reader        httpReader
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
	fmt.Printf("got some data\n")
}

// New is required to statisfy the StreamFactory inferface
func (h *HTTPStreamFactory) New(netFlow, tcpFlow gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fmt.Printf("********* starting a new stream for: %s, %s\n", netFlow, tcpFlow)
	stream := &httpStream{
		netFlow:       netFlow,
		transportFlow: tcpFlow,
	}
	// go stream.run()
	return stream
}

// func (h *httpStream) run() {
// 	buf := bufio.NewReader(&h.reader)
// 	for {
// 		req, err := http.ReadRequest(buf)
// 		if err == io.EOF {
// 			return
// 		} else if err != nil {
// 			fmt.Printf("cannot read request, %s\n", err.Error())
// 		} else {
// 			print.Request(req)
// 		}
// 	}
// }
