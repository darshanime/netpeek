package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/darshanime/netpeek/cui"
	"github.com/darshanime/netpeek/stream"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
)

var iface = flag.String("i", "eth0", "Interface to read packets from")
var useCui = flag.Bool("cui", false, "Set CUI mode")
var bpf = flag.String("bpf", "tcp port 80", "bpf program")

func main() {
	flag.Parse()
	fmt.Printf("%s\n", pcap.Version())

	handle, err := pcap.OpenLive("en0", int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic("cannot open en0 interface for sniffing")
	}
	defer handle.Close()
	err = handle.SetBPFFilter(*bpf)
	if err != nil {
		panic("incorrect bpf program")
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	streamFactory := &stream.HTTPStreamFactory{UseCui: useCui}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	if *useCui {
		go cui.InitCui()
	}
	for {
		select {
		case packet := <-packets:
			if !*useCui {
				fmt.Fprintf(os.Stdout, "#")
			}

			if packet == nil {
				return
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
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
