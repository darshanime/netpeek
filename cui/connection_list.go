package cui

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

func AddConnection(netflow, tcpflow gopacket.Flow, sno string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("conns")
		if err != nil {
			return err
		}
		if _, ok := connMap[netflow]; !ok {
			connMap[netflow] = 1
		}
		addLineToViewConns(v, sno,
			netflow.Src().String()+":"+tcpflow.Src().String(),
			netflow.Dst().String()+":"+tcpflow.Dst().String())
		return nil
	})
}

func addLineToViewConns(v *gocui.View, sno, src, dst string) {
	line := pad.Right(sno, 10, " ") + pad.Right(src, 30, " ") + pad.Right(dst, 30, " ")
	fmt.Fprintln(v, line)
}

func getConnectionListView(tcp *layers.TCP) string {
	return tcp.TransportFlow().String()
}
