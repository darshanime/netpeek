package main

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
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

	handle, _ := pcap.OpenLive("en0", int32(65535), true, -1*time.Second)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packet, err := packetSource.NextPacket()
	if err != nil {
		panic("could not get first packet")
	}
	fmt.Printf("%s\n", packet)
	defer handle.Close()
}
