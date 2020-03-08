package stream

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/darshanime/netpeek/print"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type HTTPStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	reader         tcpreader.ReaderStream
}

// New is required to statisfy the StreamFactory inferface
func (h *HTTPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	fmt.Printf("********* starting a new stream for: %s, %s\n", net, transport)
	stream := &httpStream{
		net:       net,
		transport: transport,
		reader:    tcpreader.NewReaderStream(),
	}
	go stream.run()
	return &stream.reader
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.reader)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		} else if err != nil {
			fmt.Printf("cannot read request, %s\n", err.Error())
		} else {
			print.Request(req)
		}
	}
}
