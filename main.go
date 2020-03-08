package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darshanime/netpeek/stream"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
)

func main() {
	fmt.Printf("%s\n", pcap.Version())
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic("could not FindAllDevs")
	}
	for _, device := range devices {
		fmt.Printf("device: %s\n", device.Name)
	}

	handle, err := pcap.OpenLive("en0", int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic("cannot open en0 interface for sniffing")
	}
	defer handle.Close()
	err = handle.SetBPFFilter("tcp and dst port 80")
	if err != nil {
		panic("incorrect bpf program")
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	streamFactory := &stream.HTTPStreamFactory{}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
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
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushCloseOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
