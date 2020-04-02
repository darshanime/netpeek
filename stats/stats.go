package stats

import (
	"time"
)

type PacketInfo struct {
	FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
	CaptureLength                              int
	Timestamp                                  time.Duration
	Dir                                        string
}

func (p *PacketInfo) String() string {
	return getDirection(p) + " " + getFlags(p) + p.Timestamp.String()
}

func getDirection(p *PacketInfo) string {
	if p.Dir == "client->server" {
		return "<--"
	}
	return "-->"
}

func getFlags(p *PacketInfo) string {
	return getElement(&p.FIN, "FIN") +
		getElement(&p.SYN, "SYN") +
		getElement(&p.RST, "RST") +
		getElement(&p.PSH, "PSH") +
		getElement(&p.ACK, "ACK") +
		getElement(&p.URG, "URG") +
		getElement(&p.ECE, "ECE") +
		getElement(&p.CWR, "CWR") +
		getElement(&p.NS, "NS")
}

func getElement(b *bool, s string) string {
	if *b {
		return s + " "
	}
	return "... "
}
