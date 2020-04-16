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

	handle, err := pcap.OpenLive(*iface, int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic(fmt.Sprintf("cannot open %s interface for sniffing", *iface))
	}
	defer handle.Close()
	err = handle.SetBPFFilter(getBPFProgram())
	if err != nil {
		panic("incorrect bpf program")
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
				logger.Printf("Unusable packet")
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

func getBPFProgram() string {
	flagMap := getFlagMap()
	switch flagMap {
	case "0000":
		return "tcp"
	case "0001":
		return fmt.Sprintf("tcp and (dst host %s or src host %s)", *dHost, *dHost)
	case "0010":
		return fmt.Sprintf("tcp and (src host %s or dst host %s)", *sHost, *sHost)
	case "0100":
		return fmt.Sprintf("tcp and (dst port %s or src port %s)", *dPort, *dPort)
	case "1000":
		return fmt.Sprintf("tcp and (src port %s or dst port %s)", *sPort, *sPort)
	case "1001":
		return fmt.Sprintf("tcp and ((src port %s and dst host %s) or (dst port %s and src host %s))", *sHost, *dHost, *sHost, *dHost)
	case "1010":
		return fmt.Sprintf("tcp and ((src port %s and src host %s) or (dst port %s and dst host %s))", *sPort, *sHost, *sPort, *sHost)
	case "0011":
		return fmt.Sprintf("tcp and ((src host %s and dst host %s) or (dst host %s and src host %s))", *sHost, *dHost, *sHost, *dHost)
	case "0101":
		return fmt.Sprintf("tcp and ((dst port %s and dst host %s) or (src port %s and src host %s))", *dPort, *dHost, *dPort, *dHost)
	case "0110":
		return fmt.Sprintf("tcp and ((dst port %s and src host %s) or (src port %s and dst host %s))", *dPort, *sHost, *dPort, *sHost)
	case "1100":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s) or (dst port %s and src port %s))", *sPort, *dPort, *sPort, *dPort)
	case "0111":
		return fmt.Sprintf("tcp and ((dst port %s and src host %s and dst host %s) or (src port %s and dst host %s and src host %s))", *dPort, *sHost, *dHost, *dPort, *sHost, *dHost)
	case "1011":
		return fmt.Sprintf("tcp and ((src port %s and src host %s and dst host %s) or (dst port %s and dst host %s and src host %s))", *sPort, *sHost, *dHost, *sPort, *sHost, *dHost)
	case "1101":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and dst host %s) or (dst port %s and src port %s and src host %s))", *sPort, *dPort, *dHost, *sPort, *dPort, *dHost)
	case "1110":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and src host %s) or (dst port %s and src port %s and dst host %s))", *sPort, *dPort, *sHost, *sPort, *dPort, *sHost)
	case "1111":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and src host %s and dst host %s) or (dst port %s and src port %s and dst host %s and src host %s))", *sPort, *dPort, *sHost, *dHost, *sPort, *dPort, *sHost, *dHost)
	default:
		panic("could not build bpf program")
	}
}

func getFlagMap() string {
	isPresent := func(opt *string) string {
		if *opt != "" {
			return "1"
		}
		return "0"
	}
	return isPresent(sPort) + isPresent(dPort) + isPresent(sHost) + isPresent(dHost)
}
