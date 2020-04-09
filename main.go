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

var iface = flag.String("i", "en0", "interface to read packets from")
var useCui = flag.Bool("cui", false, "set CUI mode")
var bpf = flag.String("bpf", "tcp", "bpf program")
var appPort = flag.String("port", "80", "application http port")
var quiet = flag.Bool("q", false, "quiet mode")

func main() {
	flag.Parse()
	fmt.Printf("%s\n", pcap.Version())
	fmt.Printf("iface %s\n", *iface)
	fmt.Printf("useCui %t\n", *useCui)
	fmt.Printf("bpf %s\n", *bpf)
	fmt.Printf("appPort %s\n", *appPort)
	fmt.Printf("quiet %t\n", *quiet)

	handle, err := pcap.OpenLive(*iface, int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic(fmt.Sprintf("cannot open %s interface for sniffing", *iface))
	}
	defer handle.Close()
	err = handle.SetBPFFilter(*bpf)
	if err != nil {
		panic("incorrect bpf program")
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	streamFactory := &stream.HTTPStreamFactory{UseCui: useCui, HTTPPort: appPort}
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
			if !*useCui && !*quiet {
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
