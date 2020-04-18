package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/stream"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
)

var iface = flag.String("i", "en0", "interface to read packets from")
var useCui = flag.Bool("cui", false, "set CUI mode")
var sPort = flag.String("sport", "", "source port")
var dPort = flag.String("dport", "", "destination port")
var sHost = flag.String("shost", "", "source host port")
var dHost = flag.String("dhost", "", "destination port")
var protocol = flag.String("protocol", "http", "protocol to parse for")
var verbose = flag.Bool("v", false, "verbose mode")
var bpf = flag.String("bpf", "", "custom bpf program")

func main() {
	flag.Parse()
	logger := getLogger(verbose, useCui)
	logger.Printf("pcap version: %s", pcap.Version())
	logger.Printf("iface: %s", *iface)
	logger.Printf("useCui: %t", *useCui)
	logger.Printf("sPort: %s", *sPort)
	logger.Printf("dPort: %s", *dPort)
	logger.Printf("sHost: %s", *sHost)
	logger.Printf("dHost: %s", *dHost)
	logger.Printf("protocol: %s", *protocol)
	logger.Printf("verbose: %t", *verbose)
	logger.Printf("bpf: %s", getBPFProgram())

	if *useCui && *protocol != "http" {
		panic("cui mode can only be used with http protocol")
	}

	handle, err := pcap.OpenLive(*iface, int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic(fmt.Sprintf("cannot open %s interface for sniffing", *iface))
	}
	defer handle.Close()

	if *bpf == "" {
		*bpf = getBPFProgram()
	}

	if err := handle.SetBPFFilter(*bpf); err != nil {
		panic(fmt.Sprintf("incorrect bpf program: %s", *bpf))
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	streamFactory := &stream.HTTPStreamFactory{UseCui: useCui, Protocol: protocol, Logger: logger}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	if *useCui {
		go cui.InitCui(logger)
	}
	for {
		select {
		case packet := <-packets:
			// logger adds new lines mandatorily - see golang/go/issues/16564
			if !*useCui && *verbose {
				fmt.Printf("#")
			}

			if packet == nil {
				return
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				logger.Printf("unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			c := stream.AssemblerContext{
				CaptureInfo: packet.Metadata().CaptureInfo,
			}
			assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &c)

		case <-ticker:
			assembler.FlushCloseOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
