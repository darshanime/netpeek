package stats

import (
	"html/template"
	"time"
)

type PacketInfo struct {
	FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
	CaptureLength                              int
	Timestamp                                  time.Duration
	Dir                                        template.HTML
}
